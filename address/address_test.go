package address

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	{
		path := "/"

		address := NadaAddress{}
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, path, sb.String())

		{
			address := GetAddress(path)
			_, ok := address.(NadaAddress)
			assert.True(t, ok)
		}
	}

	{
		path := "/channels"

		var address ChannelsAddress
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, path, sb.String())

		{
			address := GetAddress(path)
			_, ok := address.(ChannelsAddress)
			assert.True(t, ok)
		}
	}

	{
		path := "/topics/42"

		address := TopicsAddress{
			ChannelId: 42,
		}
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, path, sb.String())

		{
			address := GetAddress(path)
			topicsAddress, ok := address.(TopicsAddress)
			assert.True(t, ok)
			assert.Equal(
				t,
				address,
				topicsAddress,
			)
		}
	}

	{
		path := "/messages/99"

		address := MessagesAddress{
			AddressIndex: 99,
		}
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, path, sb.String())

		{
			address := GetAddress(path)
			messagesAddress, ok := address.(MessagesAddress)
			assert.True(t, ok)
			assert.Equal(
				t,
				address,
				messagesAddress,
			)
		}
	}
}
