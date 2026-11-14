package drain

import "io"

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

type Drainer interface {
	AcknowledgeReceive(size int)
	Drain(reader io.Reader) error
}
