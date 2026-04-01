package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

import (
	"go-zulip/database"
	"go-zulip/server_types"
)

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

func TestTopic(t *testing.T) {
	table := database.NewTopicTable()

	table.Put("apple")
	table.Put("apple")
	table.Put("apple")
	table.Put("banana")
	table.Put("banana")

	assert.Equal(t, table.NameFromIndex(0), "apple")
	assert.Equal(t, table.NameFromIndex(1), "banana")
}

func addTestSubs(t *testing.T, db *database.Database) {
	subs := []ServerSubscription{
		{
			Stream_id: 101,
			Name:      "engineering",
		},
		{
			Stream_id: 102,
			Name:      "design",
		},
	}

	for i, sub := range subs {
		db.AddServerSubscription(sub)
		// test idempotency
		index := db.AddServerSubscription(sub)
		assert.Equal(t, index, i)
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

func TestMessages(t *testing.T) {
	db := database.NewDatabase()
	messages := test_messages()
	addTestMessages(db, messages)
	assert.Equal(t, messages[0].Content, "message0")
	assert.Equal(t, db.UserTable.GetName(1002), "Fred Flintstone")
}

func TestChannels(t *testing.T) {
	db := database.NewDatabase()

	addTestSubs(t, db)

	assert.Equal(t, db.ChannelTable.GetName(101), "engineering")
	assert.Equal(t, db.ChannelTable.GetName(102), "design")
}
