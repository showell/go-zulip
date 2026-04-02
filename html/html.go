package html

import h "html"

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

import (
	"go-zulip/database"
)

type Database = database.Database

func ChannelsHtml(db *Database) string {
	var sb strings.Builder
	rows := db.ChannelTable.Rows

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

	p := sb.WriteString
	Itoa := strconv.Itoa

	p("<h4>")
	p(Itoa(len(rows)))
	p(" channels</h4>\n")

	for _, index := range indexes {
		row := rows[index]
		name := h.EscapeString(row.Name)
		channel_id := Itoa(row.Id)
		num_topics := Itoa(db.ChannelToAddress.Count(row.Index))

		p("<div class='channel_row'>\n<div class='channel_name'>")
		p(name)
		p("</div>\n<div><a href='/topics/")
		p(channel_id)
		p(">topic</a></div>\n<div class='channel_count'>")
		p(num_topics)
		p(" topics</div>\n</div>\n")
		p("\n")
	}

	return sb.String()
}

func TopicsHtml(db *Database, channel_index int) string {
	var sb strings.Builder

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

	p := sb.WriteString
	Itoa := strconv.Itoa

	p("<h4>")
	p(Itoa(topics_count))
	p(" topics</h4>\n")

	for _, row := range rows {
		topic_name := h.EscapeString(row.topic_name)
		address_index := Itoa(row.address_index)
		message_count := Itoa(row.message_count)

		p("<div class='topic_row'>\n<div class='topic_name'>")
		p(topic_name)
		p("</div>\n<div><a href='/topic_messages/")
		p(address_index)
		p(">topic</a></div>\n<div class='topic_count'>")
		p(message_count)
		p(" messages</div>\n</div>\n")
		p("\n")
	}

	return sb.String()
}

func MessagesHtml(db *Database, address_index int) string {
	var sb strings.Builder

	address := db.AddressTable.Rows[address_index]
	channel_index := address.ChannelIndex
	topic_index := address.TopicIndex
	topic_name := db.TopicTable.Rows[topic_index]

	message_indexes := db.AddressToMessage.GetManyIndexesInRandomOrder(channel_index)

	message_count := len(message_indexes)

	type MessageRow struct {
		sender_name string
		content     string
	}

	rows := make([]MessageRow, len(message_indexes))

	for i, message_index := range message_indexes {
		message := db.MessageTable.Rows[message_index]
		sender_name := db.UserTable.Rows[message.SenderIndex].Name
		content := message.Content

		rows[i] = MessageRow{
			sender_name: sender_name,
			content:     content,
		}
	}

	slices.SortFunc(rows, func(a, b MessageRow) int {
		return cmp.Compare(a.sender_name, b.sender_name)
	})

	p := sb.WriteString
	Itoa := strconv.Itoa

	p("<h4>")
	p(Itoa(message_count))
	p(" messages for ")
	p(h.EscapeString(topic_name))
	p("</h4>\n")

	for _, row := range rows {
		sender_name := h.EscapeString(row.sender_name)
		content := row.content // already valid HTML from server

		p("<div class='message_sender'>")
		p(sender_name)
		p("</div>\n")
		p("<div>")
		p(content)
		p("</div>\n")
	}

	return sb.String()
}
