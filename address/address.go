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

func (self ChannelsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/channels")
}

func (self TopicsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/topics/")
	w.WriteString(strconv.Itoa(self.channel_index))
}
