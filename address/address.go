package address

import "io"

type AddressType = string

type ChannelAddress struct {
	Type AddressType
}

func (self ChannelAddress) WritePath(w io.StringWriter) {
	w.WriteString("/channels")
}
