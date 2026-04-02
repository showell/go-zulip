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

	for address_index := range topics_count {
		topic_name := db.TopicTable.Rows[address_index]
		message_count := db.AddressToMessage.Count(address_index)

		rows[address_index] = TopicRow{
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
