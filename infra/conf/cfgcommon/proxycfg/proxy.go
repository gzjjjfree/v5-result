package proxycfg

import "github.com/gzjjjfree/v5-result/transport/internet"

type ProxyConfig struct {
	Tag                 string `json:"tag"`
	TransportLayerProxy bool   `json:"transportLayer"`
}

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

// Build implements Buildable.
func (v *ProxyConfig) Build() (*internet.ProxyConfig, error) {
	if v.Tag == "" {
		return nil, newError("Proxy tag is not set.")
	}
	return &internet.ProxyConfig{
		Tag:                 v.Tag,
		TransportLayerProxy: v.TransportLayerProxy,
	}, nil
}
