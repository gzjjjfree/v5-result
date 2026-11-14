//go:build !linux || (linux && !(amd64 || arm64))
// +build !linux linux,!amd64,!arm64

package gvisor

import "github.com/gzjjjfree/v5-result/app/tun/device"

func New(options device.Options) (device.Device, error) {
	return nil, newError("not supported").AtError()
}
