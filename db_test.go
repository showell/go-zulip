package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

import (
	"go-zulip/database"
	"go-zulip/html"
	"go-zulip/server_types"
)

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

func TestMessage(t *testing.T) {
	message_table := database.NewMessageTable()

	message_table.Put(database.Message{
		MessageId:    1001,
		SenderIndex:  200,
		AddressIndex: 300,
		Content:      "message 1001",
	})

	assert.Equal(
		t,
		database.MessageRow{
			Index:        0,
			MessageId:    1001,
			SenderIndex:  200,
			AddressIndex: 300,
			Content:      "message 1001",
		},
		*message_table.RowFromId(1001),
	)
}

func TestOneToMany(t *testing.T) {
	one_to_many := database.NewOneToMany()

	one_to_many.Update(0, 5)
	one_to_many.Update(3, 30)
	one_to_many.Update(2, 22)
	one_to_many.Update(0, 7)
	one_to_many.Update(0, 3)

	get := func(one_index int) []int {
		lst := one_to_many.GetManyIndexesInRandomOrder(one_index)
		slices.Sort(lst)
		return lst
	}

	assert.Equal(t, []int{3, 5, 7}, get(0))
	assert.Equal(t, []int{}, get(1))
	assert.Equal(t, []int{22}, get(2))
	assert.Equal(t, []int{30}, get(3))

	assert.Equal(t, 3, one_to_many.Count(0))
}

func TestAddress(t *testing.T) {
	table := database.NewAddressTable()

	type Row = database.AddressRow

	table.Put(Row{ChannelIndex: 0, TopicIndex: 0})
	table.Put(Row{ChannelIndex: 0, TopicIndex: 0})
	table.Put(Row{ChannelIndex: 0, TopicIndex: 0})
	table.Put(Row{ChannelIndex: 0, TopicIndex: 1})
	table.Put(Row{ChannelIndex: 0, TopicIndex: 1})
	table.Put(Row{ChannelIndex: 4, TopicIndex: 1})

	assert.Equal(
		t,
		&Row{ChannelIndex: 0, TopicIndex: 0},
		table.RowFromIndex(0),
	)
	assert.Equal(
		t,
		&Row{ChannelIndex: 0, TopicIndex: 1},
		table.RowFromIndex(1),
	)
	assert.Equal(
		t,
		&Row{ChannelIndex: 4, TopicIndex: 1},
		table.RowFromIndex(2),
	)
}

func TestTopic(t *testing.T) {
	table := database.NewTopicTable()

	table.Put("apple")
	table.Put("apple")
	table.Put("apple")
	table.Put("banana")
	table.Put("banana")

	assert.Equal(t, "apple", table.NameFromIndex(0))
	assert.Equal(t, "banana", table.NameFromIndex(1))
}

func addTestSubs(t *testing.T, db *database.Database) {
	subs := []ServerSubscription{
		{
			Stream_id: 101,
			Name:      "engineering",
		},
		{
			Stream_id: 102,
			Name:      "design & frontend",
		},
	}

	for i, sub := range subs {
		db.AddServerSubscription(sub)
		// test idempotency
		index := db.AddServerSubscription(sub)
		assert.Equal(t, i, index)
	}

}

func test_messages() []ServerMessage {
	return []ServerMessage{
		{
			Content:          "message0",
			Id:               201,
			Sender_full_name: "Foo Barson",
			Sender_id:        1001,
			Subject:          "design stuff",
			Stream_id:        102,
		},

		{
			Content:          "message1",
			Id:               202,
			Sender_full_name: "Foo Barson",
			Sender_id:        1001,
			Subject:          "design stuff",
			Stream_id:        102,
		},

		{
			Content:          "message2",
			Id:               203,
			Sender_full_name: "Fred Flintstone",
			Sender_id:        1002,
			Subject:          "feedback & other stuff",
			Stream_id:        101,
		},

		{
			Content:          "message3",
			Id:               204,
			Sender_full_name: "Fred Flintstone",
			Sender_id:        1002,
			Subject:          "another design topic",
			Stream_id:        102,
		},
	}
}

func addTestMessages(db *database.Database, messages []ServerMessage) {
	for _, message := range messages {
		db.AddServerMessage(message)
	}
}

func TestChannelsHtml(t *testing.T) {
	db := database.NewDatabase()
	addTestSubs(t, db)
	addTestMessages(db, test_messages())

	fmt.Println(html.ChannelsHtml(db))
}

func TestMessages(t *testing.T) {
	db := database.NewDatabase()
	messages := test_messages()
	addTestMessages(db, messages)
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(
		t,
		"Fred Flintstone",
		db.UserTable.GetName(1002),
	)
}

func TestChannels(t *testing.T) {
	db := database.NewDatabase()

	addTestSubs(t, db)

	assert.Equal(t, "engineering", db.ChannelTable.GetName(101))
	assert.Equal(t, "design & frontend", db.ChannelTable.GetName(102))
}

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

	for i := range 5 {
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

		fmt.Printf("%d loops\n", i)
		fmt.Printf("message_id %d\n", message_id)
	}

	return db
}

func TestPerf(t *testing.T) {
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
