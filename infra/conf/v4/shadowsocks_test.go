package v4_test

import (
	"testing"

	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/common/protocol"
	"github.com/gzjjjfree/v5-result/common/serial"
	"github.com/gzjjjfree/v5-result/infra/conf/cfgcommon"
	"github.com/gzjjjfree/v5-result/infra/conf/cfgcommon/testassist"
	v4 "github.com/gzjjjfree/v5-result/infra/conf/v4"
	"github.com/gzjjjfree/v5-result/proxy/shadowsocks"
)

func TestShadowsocksServerConfigParsing(t *testing.T) {
	creator := func() cfgcommon.Buildable {
		return new(v4.ShadowsocksServerConfig)
	}

	testassist.RunMultiTestCase(t, []testassist.TestCase{
		{
			Input: `{
				"method": "aes-256-GCM",
				"password": "v2ray-password"
			}`,
			Parser: testassist.LoadJSON(creator),
			Output: &shadowsocks.ServerConfig{
				User: &protocol.User{
					Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: shadowsocks.CipherType_AES_256_GCM,
						Password:   "v2ray-password",
					}),
				},
				Network: []net.Network{net.Network_TCP},
			},
		},
	})
}
