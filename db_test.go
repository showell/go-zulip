package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

import "go-zulip/database"
import "go-zulip/server_types"

type ServerSubscription = server_types.ServerSubscription

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

	assert.Equal(t, db.GetChannelName(101), "engineering")
	assert.Equal(t, db.GetChannelName(102), "design")
}
