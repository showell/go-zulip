package address

import (
	"io"
	"regexp"
	"strconv"
)

type Address interface {
	WritePath(w io.StringWriter)
}

type NadaAddress struct {
}

type ChannelsAddress struct {
}

type TopicsAddress struct {
	ChannelId int
}

type MessagesAddress struct {
	AddressIndex int
}

func (self NadaAddress) WritePath(w io.StringWriter) {
	w.WriteString("/")
}

func (self ChannelsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/channels")
}

func (self TopicsAddress) WritePath(w io.StringWriter) {
	w.WriteString("/topics/")
	w.WriteString(strconv.Itoa(self.ChannelId))
}

func (self MessagesAddress) WritePath(w io.StringWriter) {
	w.WriteString("/messages/")
	w.WriteString(strconv.Itoa(self.AddressIndex))
}

var topicRegex = regexp.MustCompile(`/topics/(\d+)`)
var topic_matches = topicRegex.FindStringSubmatch

var messagesRegex = regexp.MustCompile(`/messages/(\d+)`)
var messages_matches = messagesRegex.FindStringSubmatch

func GetAddress(path string) Address {
	if path == "/channels" {
		return ChannelsAddress{}
	} else if matches := topic_matches(path); matches != nil {
		channel_id, err := strconv.Atoi(matches[1])
		if err != nil {
			return NadaAddress{}
		}
		return TopicsAddress{ChannelId: channel_id}
	} else if matches := messages_matches(path); matches != nil {
		address_index, err := strconv.Atoi(matches[1])
		if err != nil {
			return NadaAddress{}
		}
		return MessagesAddress{AddressIndex: address_index}
	}

	return NadaAddress{}
}
