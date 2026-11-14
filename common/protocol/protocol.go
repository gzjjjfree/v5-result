package protocol

import (
	"errors"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

var ErrProtoNeedMoreData = errors.New("protocol matches, but need more data to complete sniffing")
