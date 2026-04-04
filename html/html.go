package html

import addr "go-zulip/address"
import h "html"

import (
	"cmp"
	"io"
	"slices"
	"strconv"
)

import (
	"go-zulip/database"
)

type Database = database.Database

func ChannelsHtml(db *Database, writer io.StringWriter) {
	rows := db.ChannelTable.Rows

	w := writer.WriteString

	indexes := make([]int, len(rows))
	for i := range len(rows) {
		indexes[i] = i
	}

	slices.SortFunc(indexes, func(a, b int) int {
		return cmp.Compare(
			rows[a].Name,
			rows[b].Name,
		)
	})

	Itoa := strconv.Itoa

	w("<h4>")
	w(Itoa(len(rows)))
	w(" channels</h4>\n")

	for _, index := range indexes {
		row := rows[index]
		name := h.EscapeString(row.Name)
		channel_id := row.Id
		num_topics := Itoa(db.ChannelToAddress.Count(row.Index))

		topics_address := addr.TopicsAddress{ChannelId: channel_id}

		w("<div class='channel_row'>\n<div class='channel_name'>")
		w(name)
		w("</div>\n<div><a href='")
		topics_address.WritePath(writer)
		w("'>topics</a></div>\n<div class='channel_count'>")
		w(num_topics)
		w(" topics</div>\n</div>\n")
		w("\n")
	}
}

func TopicsHtml(db *Database, channel_id int, writer io.StringWriter) {
	w := writer.WriteString

	channel_index := db.ChannelTable.GetOrMakeIndex(channel_id)
	channel_name := db.ChannelTable.Rows[channel_index].Name

	address_indexes := db.ChannelToAddress.GetManyIndexesInRandomOrder(channel_index)

	topics_count := len(address_indexes)

	type TopicRow struct {
		topic_name    string
		address_index int
		message_count int
	}

	rows := make([]TopicRow, len(address_indexes))

	for i, address_index := range address_indexes {
		topic_index := db.AddressTable.Rows[address_index].TopicIndex
		topic_name := db.TopicTable.Rows[topic_index]
		message_count := db.AddressToMessage.Count(address_index)

		rows[i] = TopicRow{
			topic_name:    topic_name,
			address_index: address_index,
			message_count: message_count,
		}
	}

	slices.SortFunc(rows, func(a, b TopicRow) int {
		return cmp.Compare(a.topic_name, b.topic_name)
	})

	Itoa := strconv.Itoa

	w("<h4>")
	w(Itoa(topics_count))
	w(" topics for ")
	w(h.EscapeString(channel_name))
	w("</h4>\n")

	for _, row := range rows {
		topic_name := h.EscapeString(row.topic_name)
		address_index := row.address_index
		message_count := Itoa(row.message_count)

		messages_address := addr.MessagesAddress{AddressIndex: address_index}

		w("<div class='topic_row'>\n<div class='topic_name'>")
		w(topic_name)
		w("</div>\n<div><a href='")
		messages_address.WritePath(writer)
		w("'>messages</a></div>\n<div class='topic_count'>")
		w(message_count)
		w(" messages</div>\n</div>\n")
		w("\n")
	}
}

func MessagesHtml(db *Database, address_index int, writer io.StringWriter) {
	w := writer.WriteString

	address := db.AddressTable.Rows[address_index]
	channel_index := address.ChannelIndex
	channel_name := db.ChannelTable.Rows[channel_index].Name
	topic_index := address.TopicIndex
	topic_name := db.TopicTable.Rows[topic_index]

	message_indexes := db.AddressToMessage.GetManyIndexesInRandomOrder(address_index)

	message_count := len(message_indexes)

	type MessageRow struct {
		message_id  int
		sender_name string
		content     string
	}

	rows := make([]MessageRow, len(message_indexes))

	for i, message_index := range message_indexes {
		message := db.MessageTable.Rows[message_index]
		message_id := message.MessageId
		sender_name := db.UserTable.Rows[message.SenderIndex].Name
		content := message.Content

		rows[i] = MessageRow{
			message_id:  message_id,
			sender_name: sender_name,
			content:     content,
		}
	}

	slices.SortFunc(rows, func(a, b MessageRow) int {
		return cmp.Compare(a.message_id, b.message_id)
	})

	Itoa := strconv.Itoa

	w("<h4>")
	w(Itoa(message_count))
	w(" messages for #")
	w(h.EscapeString(channel_name))
	w(" > ")
	w(h.EscapeString(topic_name))
	w("</h4>\n")

	for _, row := range rows {
		// sender_name := h.EscapeString(row.sender_name)
		sender_name := row.sender_name
		content := row.content // already valid HTML from server

		w("<div class='message_sender'>")
		w(sender_name)
		w("</div>\n")
		w("<div>")
		w(content)
		w("</div>\n")
	}

	w("<h3>")
	w("THE END: ")
	w(Itoa(message_count))
	w(" messages for ")
	w(h.EscapeString(topic_name))
	w("</h3>\n\n")
}

func Html(db *Database, path string, writer io.StringWriter) {
	address := addr.GetAddress(path)
	switch v := address.(type) {
	case addr.ChannelsAddress:
		ChannelsHtml(db, writer)
	case addr.TopicsAddress:
		TopicsHtml(db, v.ChannelId, writer)
	case addr.MessagesAddress:
		MessagesHtml(db, v.AddressIndex, writer)
	default:
		writer.WriteString("invalid path")
	}
}
