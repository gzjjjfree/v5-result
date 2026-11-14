//go:build !linux && !freebsd
// +build !linux,!freebsd

package tcp

import (
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}
