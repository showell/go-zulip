package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
)

import (
	"go-zulip/database"
	"go-zulip/html"
	"go-zulip/server_types"
	"go-zulip/zulip"
	// "strings"
)

//go:embed style.css
var styleCSS []byte

type ServerMessage = servertypes.ServerMessage
type ServerSubscription = servertypes.ServerSubscription

func buildBigDb() *database.Database {
	db := database.NewDatabase()

	nums := [20]int{17, 11, 4, 6, 14, 2, 9, 12, 1, 13, 19, 15, 5, 7, 10, 3, 8, 16, 18, 0}

	for _, n := range nums {
		channelId := 100 + n
		name := fmt.Sprintf("channel-%d", channelId)

		subscription := ServerSubscription{
			StreamId: channelId,
			Name:     name,
		}
		db.AddServerSubscription(subscription)
	}

	messageId := 0

	for range 1000 {
		for _, n := range nums {
			channelId := 100 + n

			for _, topicN := range nums {
				subject := fmt.Sprintf("topic-%d", 1000+topicN)

				messageId += 1

				if (messageId)%1_000 == 0 {
					fmt.Printf("messageId %d\n", messageId)
				}

				content := fmt.Sprintf("content for %d", messageId)

				message := ServerMessage{
					Content:        content,
					Id:             messageId,
					SenderFullName: "Foo Barson",
					SenderId:       1001,
					Subject:        subject,
					StreamId:       channelId,
				}

				db.AddServerMessage(message)
			}
		}

	}

	return db
}

type Counter struct {
	cnt int64
}

func (c *Counter) WriteString(s string) (int, error) {
	c.cnt += int64(len(s))

	return 0, nil
}

func channels() {
	db := buildBigDb()
	counter := Counter{}

	fmt.Println("Test channels html")
	for i := range 10_000_000 {
		html.ChannelsHtml(db, &counter)

		if (i+1)%500_000 == 0 {
			fmt.Println(i + 1)
		}
	}
	fmt.Println(counter.cnt)

	// sanity check
	html.ChannelsHtml(db, os.Stdout)
}

func topicsAndMessages() {
	type AddressRow = database.AddressRow

	db := buildBigDb()
	fmt.Println("Test topics html")
	counter := Counter{}

	loop := 0
	for counter.cnt < 10_000_000_000 {
		loop += 1

		if loop%20 == 0 {
			fmt.Println(loop, "(outer loop)")
		}

		for i := range 20 {
			channelId := 100 + i
			channelIndex := db.ChannelTable.GetOrMakeIndex(channelId)

			html.TopicsHtml(db, channelId, &counter)

			for j := range 20 {
				subject := fmt.Sprintf("topic-%d", 1000+j)
				topicIndex := db.TopicTable.Put(subject)

				addressRow := AddressRow{
					ChannelIndex: channelIndex,
					TopicIndex:   topicIndex,
				}
				addressIndex := db.AddressTable.Put(addressRow)

				html.MessagesHtml(db, addressIndex, &counter)
			}
		}
	}

	// sanity check
	html.MessagesHtml(db, 0, os.Stdout)
	html.TopicsHtml(db, 101, os.Stdout)
	fmt.Println(counter.cnt, "chars")
}

type StringWriterForBytes struct {
	w io.Writer
}

func (sw StringWriterForBytes) WriteString(s string) (int, error) {
	sw.w.Write([]byte(s))
	return 0, nil
}

func webServer() {
	db, err := zulip.BuildDatabase("config.json")
	if err != nil {
		fmt.Printf("Error building database: %v\n", err)
		return
	}

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(styleCSS)
	})

	http.HandleFunc("/{path...}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		stringWriter := StringWriterForBytes{w: w}
		stringWriter.WriteString("<html><head><link rel='stylesheet' href='/style.css'></head><body>\n")
		path := r.PathValue("path")
		html.Html(db, "/"+path, stringWriter)
		stringWriter.WriteString("</body></html>\n")
	})

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}

func main() {
	// channels()
	// topicsAndMessages()
	webServer()
}
