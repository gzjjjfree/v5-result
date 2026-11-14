package assembly

import (
	"context"

	"github.com/gzjjjfree/v5-result/common"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

const protocolName = "request"

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return nil, newError("request is a transport protocol.")
	}))
}
