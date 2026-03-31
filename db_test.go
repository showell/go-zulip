package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

import (
	"go-zulip/database"
	"go-zulip/server_types"
	"go-zulip/topic"
)

type ServerSubscription = server_types.ServerSubscription

func TestTopic(t *testing.T) {
	table := topic.NewTopicTable()

	table.Put("apple")
	table.Put("apple")
	table.Put("apple")
	table.Put("banana")
	table.Put("banana")

	assert.Equal(t, table.NameFromIndex(0), "apple")
	assert.Equal(t, table.NameFromIndex(1), "banana")
}

func TestGeneral(t *testing.T) {
	subs := []ServerSubscription{
		{
			StreamId: 101,
			Name:     "engineering",
		},
		{
			StreamId: 102,
			Name:     "design",
		},
	}

	db := database.NewDatabase()

	for i, sub := range subs {
		db.AddServerSubscription(sub)
		// test idempotency
		index := db.AddServerSubscription(sub)
		assert.Equal(t, index, i)
	}

	assert.Equal(t, db.ChannelTable.GetName(101), "engineering")
	assert.Equal(t, db.ChannelTable.GetName(102), "design")
}
