package entries

import "github.com/gzjjjfree/v5-result/app/subscription/specs"

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

type Converter interface {
	ConvertToAbstractServerConfig(rawConfig []byte, kindHint string) (*specs.SubscriptionServerConfig, error)
}
