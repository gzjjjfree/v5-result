package udp

import (
	"context"
	"io"

	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/buf"
	"github.com/gzjjjfree/v5-result/common/net"
)

type DispatcherI interface {
	common.Closable
	Dispatch(ctx context.Context, destination net.Destination, payload *buf.Buffer)
}

var DispatcherConnectionTerminationSignalReceiverMark = "DispatcherConnectionTerminationSignalReceiverMark"

type DispatcherConnectionTerminationSignalReceiver interface {
	io.Closer
}
