package dtls

import (
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/transport/internet"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

const protocolName = "dtls"

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
