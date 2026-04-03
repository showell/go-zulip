package main

import (
	"fmt"
	"os"
)

import (
	"go-zulip/database"
	"go-zulip/html"
	"go-zulip/server_types"
	// "strings"
)

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

func build_big_db() *database.Database {
	db := database.NewDatabase()

	nums := [20]int{17, 11, 4, 6, 14, 2, 9, 12, 1, 13, 19, 15, 5, 7, 10, 3, 8, 16, 18, 0}

	for _, n := range nums {
		channel_id := 100 + n
		name := fmt.Sprintf("channel-%d", channel_id)

		subscription := ServerSubscription{
			Stream_id: channel_id,
			Name:      name,
		}
		db.AddServerSubscription(subscription)
	}

	message_id := 0

	for range 100 {
		for _, n := range nums {
			channel_id := 100 + n

			for _, topic_n := range nums {
				subject := fmt.Sprintf("topic-%d", 1000+topic_n)

				message_id += 1

				if (message_id)%1_000 == 0 {
					fmt.Printf("message_id %d\n", message_id)
				}

				content := fmt.Sprintf("content for %d", message_id)

				message := ServerMessage{
					Content:          content,
					Id:               message_id,
					Sender_full_name: "Foo Barson",
					Sender_id:        1001,
					Subject:          subject,
					Stream_id:        channel_id,
				}

				db.AddServerMessage(message)
			}
		}

	}

	return db
}

type Counter struct {
	cnt int
}

func (c *Counter) WriteString(s string) (int, error) {
	c.cnt += len(s)

	return 0, nil
}

func channels() {
	db := build_big_db()
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

func topics_and_messages() {
	type AddressRow = database.AddressRow

	db := build_big_db()
	fmt.Println("Test topics html")
	cnt := 0

	for range 100 {
		for i := range 20 {
			channel_id := 100 + i
			channel_index := db.ChannelTable.GetOrMakeIndex(channel_id)
			// topics_html := html.TopicsHtml(db, channel_index)

			for j := range 20 {
				subject := fmt.Sprintf("topic-%d", 1000+j)
				topic_index := db.TopicTable.Put(subject)

				address_row := AddressRow{
					ChannelIndex: channel_index,
					TopicIndex:   topic_index,
				}
				address_index := db.AddressTable.Put(address_row)

				messages_html := html.MessagesHtml(db, address_index)

				cnt += 1

				if cnt%100 == 0 {
					fmt.Println(cnt)
				}

				if cnt%999 == 0 {
					fmt.Println(messages_html)
				}
			}
		}
	}
}

func main() {
	channels()
	// topics_and_messages()
}
