package address

import "io"

type ChannelAddress struct {
}

func (self ChannelAddress) WritePath(w io.StringWriter) {
	w.WriteString("/channels")
}
