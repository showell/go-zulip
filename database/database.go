package database

import "go-zulip/server_types"

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

type Database struct {
	ChannelTable *IdNameTable
	UserTable    *IdNameTable
}

func NewDatabase() *Database {
	return &Database{
		ChannelTable: NewIdNameTable(),
		UserTable:    NewIdNameTable(),
	}
}

func (db *Database) AddServerSubscription(sub ServerSubscription) int {
	id := sub.Stream_id
	name := sub.Name
	return db.ChannelTable.Put(IdName{
		Id:   id,
		Name: name,
	})
}

func (db *Database) AddServerMessage(message ServerMessage) {
	db.UserTable.Put(IdName{
		Id:   message.Sender_id,
		Name: message.Sender_full_name,
	})
}
