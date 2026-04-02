package html

import (
	"go-zulip/database"
)

func ChannelsHtml(db Database) string {
	rows := db.ChannelTable.Rows
}
