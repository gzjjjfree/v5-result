package v2jsonpb

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	core "github.com/gzjjjfree/v5-result"
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/buf"
	"github.com/gzjjjfree/v5-result/common/cmdarg"
	"github.com/gzjjjfree/v5-result/common/serial"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

func loadV2JsonPb(data []byte) (*core.Config, error) {
	coreconf := &core.Config{}
	jsonpbloader := &protojson.UnmarshalOptions{Resolver: anyresolverv2{serial.GetResolver()}, AllowPartial: true}
	err := jsonpbloader.Unmarshal(data, &V2JsonProtobufFollower{coreconf.ProtoReflect()})
	if err != nil {
		return nil, err
	}
	return coreconf, nil
}

func dumpV2JsonPb(config proto.Message) ([]byte, error) {
	jsonpbdumper := &protojson.MarshalOptions{Resolver: anyresolverv2{serial.GetResolver()}, AllowPartial: true}
	bytew, err := jsonpbdumper.Marshal(&V2JsonProtobufFollower{config.ProtoReflect()})
	if err != nil {
		return nil, err
	}
	return bytew, nil
}

func DumpV2JsonPb(config proto.Message) ([]byte, error) {
	return dumpV2JsonPb(config)
}

const FormatProtobufV2JSONPB = "v2jsonpb"

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      []string{FormatProtobufV2JSONPB},
		Extension: []string{".v2pb.json", ".v2pbjson"},
		Loader: func(input interface{}) (*core.Config, error) {
			switch v := input.(type) {
			case string:
				r, err := cmdarg.LoadArg(v)
				if err != nil {
					return nil, err
				}
				data, err := buf.ReadAllToBytes(r)
				if err != nil {
					return nil, err
				}
				return loadV2JsonPb(data)
			case []byte:
				return loadV2JsonPb(v)
			case io.Reader:
				data, err := buf.ReadAllToBytes(v)
				if err != nil {
					return nil, err
				}
				return loadV2JsonPb(data)
			default:
				return nil, newError("unknown type")
			}
		},
	}))
}
