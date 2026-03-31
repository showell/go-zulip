package database

import "go-zulip/server_types"

type ServerSubscription = server_types.ServerSubscription

type Database struct {
	ChannelTable *IdNameTable
}

func NewDatabase() *Database {
	return &Database{
		ChannelTable: NewIdNameTable(),
	}
}

func (db *Database) AddServerSubscription(sub ServerSubscription) int {
	id := sub.StreamId
	name := sub.Name
	return db.ChannelTable.Put(IdName{
		Id:   id,
		Name: name,
	})
}
