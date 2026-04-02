package database

import "go-zulip/server_types"

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

type Database struct {
	AddressTable *AddressTable
	ChannelTable *IdNameTable
	TopicTable   *TopicTable
	UserTable    *IdNameTable

	// OneToMany objects are for speed
	AddressToMessage *OneToMany
	ChannelToAddress *OneToMany
}

func NewDatabase() *Database {
	return &Database{
		AddressTable: NewAddressTable(),
		ChannelTable: NewIdNameTable(),
		TopicTable:   NewTopicTable(),
		UserTable:    NewIdNameTable(),

		// OneToMany objects are for speed
		AddressToMessage: NewOneToMany(),
		ChannelToAddress: NewOneToMany(),
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
