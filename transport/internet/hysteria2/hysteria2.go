package hysteria2

import (
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/transport/internet"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

const (
	protocolName = "hysteria2"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
