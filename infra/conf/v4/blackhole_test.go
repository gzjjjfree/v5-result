package v4_test

import (
	"testing"

	"github.com/gzjjjfree/v5-result/common/serial"
	"github.com/gzjjjfree/v5-result/infra/conf/cfgcommon"
	"github.com/gzjjjfree/v5-result/infra/conf/cfgcommon/testassist"
	v4 "github.com/gzjjjfree/v5-result/infra/conf/v4"
	"github.com/gzjjjfree/v5-result/proxy/blackhole"
)

func TestHTTPResponseJSON(t *testing.T) {
	creator := func() cfgcommon.Buildable {
		return new(v4.BlackholeConfig)
	}

	testassist.RunMultiTestCase(t, []testassist.TestCase{
		{
			Input: `{
				"response": {
					"type": "http"
				}
			}`,
			Parser: testassist.LoadJSON(creator),
			Output: &blackhole.Config{
				Response: serial.ToTypedMessage(&blackhole.HTTPResponse{}),
			},
		},
		{
			Input:  `{}`,
			Parser: testassist.LoadJSON(creator),
			Output: &blackhole.Config{},
		},
	})
}
