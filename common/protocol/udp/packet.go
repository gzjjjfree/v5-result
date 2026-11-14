package udp

import (
	"github.com/gzjjjfree/v5-result/common/buf"
	"github.com/gzjjjfree/v5-result/common/net"
)

// Packet is a UDP packet together with its source and destination address.
type Packet struct {
	Payload *buf.Buffer
	Source  net.Destination
	Target  net.Destination
}
