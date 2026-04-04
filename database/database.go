package database

import (
	"go-zulip/server_types"
	"strings"
)

type ServerMessage = servertypes.ServerMessage
type ServerSubscription = servertypes.ServerSubscription

type Database struct {
	AddressTable AddressTable
	ChannelTable IdNameTable
	MessageTable MessageTable
	TopicTable   TopicTable
	UserTable    IdNameTable

	// OneToMany objects are for speed
	AddressToMessage OneToMany
	ChannelToAddress OneToMany
}

func NewDatabase() *Database {
	return &Database{
		AddressTable: *NewAddressTable(),
		ChannelTable: *NewIdNameTable(),
		MessageTable: *NewMessageTable(),
		TopicTable:   *NewTopicTable(),
		UserTable:    *NewIdNameTable(),

		// OneToMany objects are for speed
		AddressToMessage: *NewOneToMany(),
		ChannelToAddress: *NewOneToMany(),
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

func (db *Database) AddServerMessage(serverMessage ServerMessage) {
	channelId := serverMessage.StreamId
	content := serverMessage.Content
	messageId := serverMessage.Id
	senderId := serverMessage.SenderId
	senderName := serverMessage.SenderFullName
	topicName := serverMessage.Subject

	// Sender
	senderIndex := db.UserTable.Put(IdName{
		Id:   senderId,
		Name: senderName,
	})

	// Address
	channelIndex := db.ChannelTable.GetOrMakeIndex(channelId)
	topicIndex := db.TopicTable.Put(topicName)
	addressIndex := db.AddressTable.Put(AddressRow{
		ChannelIndex: channelIndex,
		TopicIndex:   topicIndex,
	})

	// Message
	message := Message{
		AddressIndex: addressIndex,
		Content:      strings.Clone(content),
		MessageId:    messageId,
		SenderIndex:  senderIndex,
	}
	messageIndex := db.MessageTable.Put(message)

	// OneToMany optimizations
	db.AddressToMessage.Update(addressIndex, messageIndex)
	db.ChannelToAddress.Update(channelIndex, addressIndex)
}
