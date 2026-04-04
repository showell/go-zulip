package zulip

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

import (
	"go-zulip/database"
	servertypes "go-zulip/server_types"
)

const initialBatchSize = 1000
const backfillBatchSize = 5000
const maxMessages = 50_000

type Config struct {
	Email   string `json:"email"`
	ApiKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	return config, nil
}

type Client struct {
	config Config
	http   *http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
		http:   &http.Client{},
	}
}

func (c *Client) authHeader() string {
	creds := c.config.Email + ":" + c.config.ApiKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func (c *Client) get(path string, params url.Values) ([]byte, error) {
	u, err := url.Parse(c.config.BaseURL + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type subscriptionsResponse struct {
	Subscriptions []servertypes.ServerSubscription `json:"subscriptions"`
}

func (c *Client) fetchSubscriptions() ([]servertypes.ServerSubscription, error) {
	body, err := c.get("api/v1/users/me/subscriptions", url.Values{})
	if err != nil {
		return nil, err
	}
	var result subscriptionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Subscriptions, nil
}

type messagesResponse struct {
	Messages    []servertypes.ServerMessage `json:"messages"`
	FoundOldest bool                        `json:"found_oldest"`
}

func (c *Client) fetchMessages(anchor string, numBefore int) (messagesResponse, error) {
	params := url.Values{}
	params.Set("narrow", "[]")
	params.Set("num_before", strconv.Itoa(numBefore))
	params.Set("anchor", anchor)

	body, err := c.get("api/v1/messages", params)
	if err != nil {
		return messagesResponse{}, err
	}
	var result messagesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return messagesResponse{}, err
	}
	return result, nil
}

func addStreamMessages(db *database.Database, messages []servertypes.ServerMessage) {
	for _, msg := range messages {
		if msg.Type == "stream" {
			db.AddServerMessage(msg)
		}
	}
}

func BuildDatabase(configPath string) (*database.Database, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	client := NewClient(config)
	db := database.NewDatabase()

	// Fetch subscriptions
	subs, err := client.fetchSubscriptions()
	if err != nil {
		return nil, fmt.Errorf("fetching subscriptions: %w", err)
	}
	for _, sub := range subs {
		db.AddServerSubscription(sub)
	}
	fmt.Printf("%d channels fetched\n", len(subs))

	// Fetch initial messages
	resp, err := client.fetchMessages("newest", initialBatchSize)
	if err != nil {
		return nil, fmt.Errorf("fetching initial messages: %w", err)
	}
	addStreamMessages(db, resp.Messages)
	fmt.Printf("%d messages fetched\n", len(db.MessageTable.Rows))

	if len(resp.Messages) == 0 {
		return db, nil
	}

	foundOldest := resp.FoundOldest
	oldestId := resp.Messages[0].Id

	// Backfill older messages
	for !foundOldest && len(db.MessageTable.Rows) < maxMessages {
		numBefore := min(maxMessages-len(db.MessageTable.Rows), backfillBatchSize)

		time.Sleep(500 * time.Millisecond)

		resp, err = client.fetchMessages(strconv.Itoa(oldestId), numBefore)
		if err != nil {
			return nil, fmt.Errorf("backfilling messages: %w", err)
		}

		addStreamMessages(db, resp.Messages)
		fmt.Printf("%d messages fetched (backfill)\n", len(db.MessageTable.Rows))

		if len(resp.Messages) == 0 {
			break
		}

		foundOldest = resp.FoundOldest
		oldestId = resp.Messages[0].Id
	}

	return db, nil
}
