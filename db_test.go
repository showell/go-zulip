package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
		Row{ChannelIndex: 0, TopicIndex: 0},
		table.Rows[0],
	)
	assert.Equal(
		t,
		Row{ChannelIndex: 0, TopicIndex: 1},
		table.Rows[1],
	)
	assert.Equal(
		t,
		Row{ChannelIndex: 4, TopicIndex: 1},
		table.Rows[2],
	)
}

func TestTopic(t *testing.T) {
	table := database.NewTopicTable()

	table.Put("apple")
	table.Put("apple")
	table.Put("apple")
	table.Put("banana")
	table.Put("banana")

	assert.Equal(t, "apple", table.Rows[0])
	assert.Equal(t, "banana", table.Rows[1])
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

	writer := bufio.NewWriter(os.Stdout)
	html.Html(db, "/channels", writer)
	writer.Flush()
}

func TestTopicsHtml(t *testing.T) {
	db := database.NewDatabase()
	addTestSubs(t, db)
	addTestMessages(db, test_messages())

	writer := bufio.NewWriter(os.Stdout)
	for _, row := range db.ChannelTable.Rows {
		channel_id := row.Id
		path := fmt.Sprintf("/topics/%d", channel_id)
		html.Html(db, path, writer)
	}
	writer.Flush()
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
