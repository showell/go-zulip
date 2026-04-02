package main

import (
	"fmt"
	"go-zulip/database"
	"go-zulip/html"
	"go-zulip/server_types"
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

	for range 25_000 {
		for n := range nums {
			channel_id := 100 + n

			for topic_n := range nums {
				subject := fmt.Sprintf("topic-%d", 1000+topic_n)

				message_id += 1

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

		if (message_id)%100_000 == 0 {
			fmt.Printf("message_id %d\n", message_id)
		}
	}

	return db
}

func channels() {
	db := build_big_db()
	fmt.Println("Test channels html")
	for i := range 50_000_001 {
		s := html.ChannelsHtml(db)
		if i%1_000_000 == 0 {
			fmt.Println(i)
		}

		if i%10_000_000 == 0 {
			fmt.Println(s)
		}
	}
}

func topics() {
	db := build_big_db()
	fmt.Println("Test topics html")
	cnt := 0
	for range 2_500_000 {
		for channel_index := range 20 {
			s := html.TopicsHtml(db, channel_index)
			cnt += 1

			if cnt%100_000 == 0 {
				fmt.Println(cnt)
			}

			if cnt%999_999 == 0 {
				fmt.Println(s)
			}
		}
	}
}

func main() {
	topics()
}
