package httpupgrade

import (
	"context"

	"github.com/gzjjjfree/v5-result/common"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

const protocolName = "httpupgrade"

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return nil, newError("httpupgrade is a transport protocol.")
	}))
}
