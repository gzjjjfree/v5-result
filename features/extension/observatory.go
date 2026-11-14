package extension

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/gzjjjfree/v5-result/features"
)

type Observatory interface {
	features.Feature
	GetObservation(ctx context.Context) (proto.Message, error)
}

func ObservatoryType() interface{} {
	return (*Observatory)(nil)
}
