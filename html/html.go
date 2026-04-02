package html

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

	for _, row := range db.ChannelTable.Rows {
		sb.WriteString(strconv.Itoa(row.Id))
		sb.WriteString("\n")
	}

	return sb.String()
}
