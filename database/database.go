package database

import (
	"go-zulip/server_types"
	"strings"
)

type ServerMessage = server_types.ServerMessage
type ServerSubscription = server_types.ServerSubscription

type Database struct {
	AddressTable *AddressTable
	ChannelTable *IdNameTable
	MessageTable *MessageTable
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
		MessageTable: NewMessageTable(),
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

func (db *Database) AddServerMessage(server_message ServerMessage) {
	channel_id := server_message.Stream_id
	content := server_message.Content
	message_id := server_message.Id
	sender_id := server_message.Sender_id
	sender_name := server_message.Sender_full_name
	topic_name := server_message.Subject

	// Sender
	sender_index := db.UserTable.Put(IdName{
		Id:   sender_id,
		Name: sender_name,
	})

	// Address
	channel_index := db.ChannelTable.GetOrMakeIndex(channel_id)
	topic_index := db.TopicTable.Put(topic_name)
	address_index := db.AddressTable.Put(AddressRow{
		ChannelIndex: channel_index,
		TopicIndex:   topic_index,
	})

	// Message
	message := Message{
		AddressIndex: address_index,
		Content:      strings.Clone(content),
		MessageId:    message_id,
		SenderIndex:  sender_index,
	}
	message_index := db.MessageTable.Put(message)

	// OneToMany optimizations
	db.AddressToMessage.Update(address_index, message_index)
	db.ChannelToAddress.Update(channel_index, address_index)
}
