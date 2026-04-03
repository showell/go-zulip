package address

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	var channel_address ChannelAddress
	var sb strings.Builder

	channel_address.WritePath(&sb)

	assert.Equal(t, "/channels", sb.String())
}
