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

	p := sb.WriteString

	for _, row := range db.ChannelTable.Rows {
		p(strconv.Itoa(row.Id))
		p(row.Name)
		p("\n")
	}

	return sb.String()
}
