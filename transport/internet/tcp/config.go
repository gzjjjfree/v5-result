package tcp

import (
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/transport/internet"
)

const protocolName = "tcp"

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
