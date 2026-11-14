package restfulapi

import (
	"context"

	"github.com/gzjjjfree/v5-result/common"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return newRestfulService(ctx, config.(*Config))
	}))
}
