package server

import (
	"context"

	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/transport/internet"
	"github.com/gzjjjfree/v5-result/transport/internet/transportcommon"
)

func ListenTLSMirror(ctx context.Context, address net.Address, port net.Port,
	streamSettings *internet.MemoryStreamConfig, handler internet.ConnHandler,
) (internet.Listener, error) {
	tlsMirrorSettings := streamSettings.ProtocolSettings.(*Config)
	listener, err := transportcommon.ListenWithSecuritySettings(ctx, address, port, streamSettings)
	if err != nil {
		return nil, newError("failed to listen TLS mirror").Base(err)
	}

	tlsMirrorServer, err := NewServer(ctx, listener, tlsMirrorSettings, handler)
	if err != nil {
		return nil, newError("failed to create TLS mirror server").Base(err)
	}

	go tlsMirrorServer.accepts()

	return tlsMirrorServer, nil
}

const protocolName = "tlsmirror"

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, ListenTLSMirror))
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
