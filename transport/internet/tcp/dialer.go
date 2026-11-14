package tcp

import (
	"context"

	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/common/serial"
	"github.com/gzjjjfree/v5-result/common/session"
	"github.com/gzjjjfree/v5-result/transport/internet"
	"github.com/gzjjjfree/v5-result/transport/internet/security"
)

// Dial dials a new TCP connection to the given destination.
func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	newError("dialing TCP to ", dest).WriteToLog(session.ExportIDToError(ctx))
	conn, err := internet.DialSystem(ctx, dest, streamSettings.SocketSettings)
	if err != nil {
		return nil, err
	}

	securityEngine, err := security.CreateSecurityEngineFromSettings(ctx, streamSettings)
	if err != nil {
		return nil, newError("unable to create security engine").Base(err)
	}

	if securityEngine != nil {
		conn, err = securityEngine.Client(conn, security.OptionWithDestination{Dest: dest})
		if err != nil {
			return nil, newError("unable to create security protocol client from security engine").Base(err)
		}
	}

	tcpSettings := streamSettings.ProtocolSettings.(*Config)
	if tcpSettings.HeaderSettings != nil {
		headerConfig, err := serial.GetInstanceOf(tcpSettings.HeaderSettings)
		if err != nil {
			return nil, newError("failed to get header settings").Base(err).AtError()
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, newError("failed to create header authenticator").Base(err).AtError()
		}
		conn = auth.Client(conn)
	}
	return internet.Connection(conn), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
