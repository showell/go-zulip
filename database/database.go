package database

import "go-zulip/server_types"
import "go-zulip/channel"

type ServerSubscription = server_types.ServerSubscription

type Database struct {
	channel_table *channel.ChannelTable
}

func NewDatabase() *Database {
	return &Database{
		channel_table: channel.NewChannelTable(),
	}
}

func (db *Database) AddServerSubscription(sub ServerSubscription) int {
	id := sub.StreamId
	name := sub.Name
	return db.channel_table.Put(channel.Channel{
		Id:   id,
		Name: name,
	})
}

func (db Database) GetChannelName(channel_id int) string {
	row := db.channel_table.RowFromId(channel_id)
	if row == nil {
		return "unknown"
	}
	return row.Name
}
