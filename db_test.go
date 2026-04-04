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

type ServerMessage = servertypes.ServerMessage
type ServerSubscription = servertypes.ServerSubscription

func TestMessage(t *testing.T) {
	messageTable := database.NewMessageTable()

	messageTable.Put(database.Message{
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
		*messageTable.RowFromId(1001),
	)
}

func TestOneToMany(t *testing.T) {
	oneToMany := database.NewOneToMany()

	oneToMany.Update(0, 5)
	oneToMany.Update(3, 30)
	oneToMany.Update(2, 22)
	oneToMany.Update(0, 7)
	oneToMany.Update(0, 3)

	get := func(oneIndex int) []int {
		lst := oneToMany.GetManyIndexesInRandomOrder(oneIndex)
		slices.Sort(lst)
		return lst
	}

	assert.Equal(t, []int{3, 5, 7}, get(0))
	assert.Equal(t, []int{}, get(1))
	assert.Equal(t, []int{22}, get(2))
	assert.Equal(t, []int{30}, get(3))

	assert.Equal(t, 3, oneToMany.Count(0))
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
			StreamId: 101,
			Name:     "engineering",
		},
		{
			StreamId: 102,
			Name:     "design & frontend",
		},
	}

	for i, sub := range subs {
		db.AddServerSubscription(sub)
		// test idempotency
		index := db.AddServerSubscription(sub)
		assert.Equal(t, i, index)
	}

}

func testMessages() []ServerMessage {
	return []ServerMessage{
		{
			Content:        "message0",
			Id:             201,
			SenderFullName: "Foo Barson",
			SenderId:       1001,
			Subject:        "design stuff",
			StreamId:       102,
		},

		{
			Content:        "message1",
			Id:             202,
			SenderFullName: "Foo Barson",
			SenderId:       1001,
			Subject:        "design stuff",
			StreamId:       102,
		},

		{
			Content:        "message2",
			Id:             203,
			SenderFullName: "Fred Flintstone",
			SenderId:       1002,
			Subject:        "feedback & other stuff",
			StreamId:       101,
		},

		{
			Content:        "message3",
			Id:             204,
			SenderFullName: "Fred Flintstone",
			SenderId:       1002,
			Subject:        "another design topic",
			StreamId:       102,
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
	addTestMessages(db, testMessages())

	writer := bufio.NewWriter(os.Stdout)
	html.Html(db, "/channels", writer)
	writer.Flush()
}

func TestTopicsHtml(t *testing.T) {
	db := database.NewDatabase()
	addTestSubs(t, db)
	addTestMessages(db, testMessages())

	writer := bufio.NewWriter(os.Stdout)
	for _, row := range db.ChannelTable.Rows {
		channelId := row.Id
		path := fmt.Sprintf("/topics/%d", channelId)
		html.Html(db, path, writer)
	}
	writer.Flush()
}

func TestMessagesHtml(t *testing.T) {
	db := database.NewDatabase()
	addTestSubs(t, db)
	addTestMessages(db, testMessages())

	writer := bufio.NewWriter(os.Stdout)
	for addressIndex := range db.AddressTable.Rows {
		path := fmt.Sprintf("/messages/%d", addressIndex)
		html.Html(db, path, writer)
	}
	writer.Flush()
}

func TestMessages(t *testing.T) {
	db := database.NewDatabase()
	messages := testMessages()
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
