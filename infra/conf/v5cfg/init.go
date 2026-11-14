package v5cfg

import (
	"bytes"
	"io"

	core "github.com/gzjjjfree/v5-result"
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/buf"
	"github.com/gzjjjfree/v5-result/common/cmdarg"
	"github.com/gzjjjfree/v5-result/infra/conf/json"
)

const jsonV5 = "jsonv5"

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      []string{jsonV5},
		Extension: []string{".v5.json", ".v5.jsonc"},
		Loader: func(input interface{}) (*core.Config, error) {
			switch v := input.(type) {
			case string:
				r, err := cmdarg.LoadArg(v)
				if err != nil {
					return nil, err
				}
				data, err := buf.ReadAllToBytes(&json.Reader{
					Reader: r,
				})
				if err != nil {
					return nil, err
				}
				return loadJSONConfig(data)
			case []byte:
				r := &json.Reader{
					Reader: bytes.NewReader(v),
				}
				data, err := buf.ReadAllToBytes(r)
				if err != nil {
					return nil, err
				}
				return loadJSONConfig(data)
			case io.Reader:
				data, err := buf.ReadAllToBytes(&json.Reader{
					Reader: v,
				})
				if err != nil {
					return nil, err
				}
				return loadJSONConfig(data)
			default:
				return nil, newError("unknown type")
			}
		},
	}))
}
