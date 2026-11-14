package command_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"

	core "github.com/gzjjjfree/v5-result"
	"github.com/gzjjjfree/v5-result/app/dispatcher"
	"github.com/gzjjjfree/v5-result/app/log"
	. "github.com/gzjjjfree/v5-result/app/log/command"
	"github.com/gzjjjfree/v5-result/app/proxyman"
	_ "github.com/gzjjjfree/v5-result/app/proxyman/inbound"
	_ "github.com/gzjjjfree/v5-result/app/proxyman/outbound"
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/serial"
)

func TestLoggerRestart(t *testing.T) {
	v, err := core.New(&core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&log.Config{}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	})
	common.Must(err)
	common.Must(v.Start())

	server := &LoggerServer{
		V: v,
	}
	common.Must2(server.RestartLogger(context.Background(), &RestartLoggerRequest{}))
}
