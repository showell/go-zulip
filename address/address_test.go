package address

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	{
		var channels_address ChannelsAddress
		var sb strings.Builder
		channels_address.WritePath(&sb)
		assert.Equal(t, "/channels", sb.String())
	}

	{
		topics_address := TopicsAddress{
			channel_index: 42,
		}
		var sb strings.Builder
		topics_address.WritePath(&sb)
		assert.Equal(t, "/topics/42", sb.String())
	}
}
