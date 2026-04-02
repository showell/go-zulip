package html

import h "html"

import (
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

	p := sb.WriteString
	Itoa := strconv.Itoa

	p("<h4>")
	p(Itoa(len(rows)))
	p(" channels</h4>\n")

	for _, row := range rows {
		name := h.EscapeString(row.Name)
		channel_id := Itoa(row.Id)
		num_topics := Itoa(db.ChannelToAddress.Count(row.Index))

		p("<div class='channel_row'>\n<div class='channel_name'>")
		p(name)
		p("</div>\n<div><a href='/topics/")
		p(channel_id)
		p(">topic</a></div>\n<div class='channel_count'>")
		p(num_topics)
		p("</div>\n</div>\n")
		p("\n")
	}

	return sb.String()
}
