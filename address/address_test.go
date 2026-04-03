package address

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	{
		var address ChannelsAddress
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, "/channels", sb.String())
	}

	{
		address := TopicsAddress{
			channel_index: 42,
		}
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, "/topics/42", sb.String())
	}

	{
		address := MessagesAddress{
			address_index: 99,
		}
		var sb strings.Builder
		address.WritePath(&sb)
		assert.Equal(t, "/messages/99", sb.String())
	}
}
