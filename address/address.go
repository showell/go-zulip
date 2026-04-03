package address

import (
	"io"
	"strconv"
)

type ChannelsAddress struct {
}

type TopicsAddress struct {
	channel_index int
}

type MessagesAddress struct {
	address_index int
}

func (self ChannelsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/channels")
}

func (self TopicsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/topics/")
	w.WriteString(strconv.Itoa(self.channel_index))
}

func (self MessagesAddress) WritePath(w io.StringWriter) {
	w.WriteString("/messages/")
	w.WriteString(strconv.Itoa(self.address_index))
}
