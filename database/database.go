package database

import "zulip-go/server_types"

type ServerSubscription = server_types.ServerSubscription

type Database struct {
    channel_map map[int]string
}

func NewDatabase() *Database {
    return &Database{
        channel_map: make(map[int]string),
    }
}

func (db *Database) AddServerSubscription(sub ServerSubscription) {
    channel_id := sub.StreamId
    name := sub.Name
    db.channel_map[channel_id] = name
}

func (db Database) GetChannelName(channel_id int) string {
    return db.channel_map[channel_id]
}

