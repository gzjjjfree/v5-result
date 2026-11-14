package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Mandatory features. Can't remove unless there are replacements.
	_ "github.com/gzjjjfree/v5-result/app/dispatcher"
	_ "github.com/gzjjjfree/v5-result/app/proxyman/inbound"
	_ "github.com/gzjjjfree/v5-result/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/gzjjjfree/v5-result/app/commander"
	_ "github.com/gzjjjfree/v5-result/app/log/command"
	_ "github.com/gzjjjfree/v5-result/app/proxyman/command"
	_ "github.com/gzjjjfree/v5-result/app/stats/command"

	// Developer preview services
	_ "github.com/gzjjjfree/v5-result/app/instman/command"
	_ "github.com/gzjjjfree/v5-result/app/observatory/command"

	// Other optional features.
	_ "github.com/gzjjjfree/v5-result/app/dns"
	_ "github.com/gzjjjfree/v5-result/app/dns/fakedns"
	_ "github.com/gzjjjfree/v5-result/app/log"
	_ "github.com/gzjjjfree/v5-result/app/policy"
	_ "github.com/gzjjjfree/v5-result/app/reverse"
	_ "github.com/gzjjjfree/v5-result/app/router"
	_ "github.com/gzjjjfree/v5-result/app/stats"

	// Fix dependency cycle caused by core import in internet package
	_ "github.com/gzjjjfree/v5-result/transport/internet/tagged/taggedimpl"

	// Developer preview features
	_ "github.com/gzjjjfree/v5-result/app/commander/webcommander"
	_ "github.com/gzjjjfree/v5-result/app/instman"
	_ "github.com/gzjjjfree/v5-result/app/observatory"
	_ "github.com/gzjjjfree/v5-result/app/persistentstorage/filesystemstorage"
	_ "github.com/gzjjjfree/v5-result/app/tun"

	// Inbound and outbound proxies.
	_ "github.com/gzjjjfree/v5-result/proxy/blackhole"
	_ "github.com/gzjjjfree/v5-result/proxy/dns"
	_ "github.com/gzjjjfree/v5-result/proxy/dokodemo"
	_ "github.com/gzjjjfree/v5-result/proxy/freedom"
	_ "github.com/gzjjjfree/v5-result/proxy/http"
	_ "github.com/gzjjjfree/v5-result/proxy/shadowsocks"
	_ "github.com/gzjjjfree/v5-result/proxy/socks"
	_ "github.com/gzjjjfree/v5-result/proxy/trojan"
	_ "github.com/gzjjjfree/v5-result/proxy/vless/inbound"
	_ "github.com/gzjjjfree/v5-result/proxy/vless/outbound"
	_ "github.com/gzjjjfree/v5-result/proxy/vmess/inbound"
	_ "github.com/gzjjjfree/v5-result/proxy/vmess/outbound"

	// Developer preview proxies
	_ "github.com/gzjjjfree/v5-result/proxy/vlite/inbound"
	_ "github.com/gzjjjfree/v5-result/proxy/vlite/outbound"

	_ "github.com/gzjjjfree/v5-result/proxy/hysteria2"
	_ "github.com/gzjjjfree/v5-result/proxy/shadowsocks2022"

	// Transports
	_ "github.com/gzjjjfree/v5-result/transport/internet/domainsocket"
	_ "github.com/gzjjjfree/v5-result/transport/internet/grpc"
	_ "github.com/gzjjjfree/v5-result/transport/internet/http"
	_ "github.com/gzjjjfree/v5-result/transport/internet/kcp"
	_ "github.com/gzjjjfree/v5-result/transport/internet/quic"
	_ "github.com/gzjjjfree/v5-result/transport/internet/tcp"
	_ "github.com/gzjjjfree/v5-result/transport/internet/tls"
	_ "github.com/gzjjjfree/v5-result/transport/internet/tls/utls"
	_ "github.com/gzjjjfree/v5-result/transport/internet/udp"
	_ "github.com/gzjjjfree/v5-result/transport/internet/websocket"

	// Developer preview transports
	_ "github.com/gzjjjfree/v5-result/transport/internet/request/assembly"

	_ "github.com/gzjjjfree/v5-result/transport/internet/request/assembler/simple"
	_ "github.com/gzjjjfree/v5-result/transport/internet/request/roundtripper/httprt"

	_ "github.com/gzjjjfree/v5-result/transport/internet/request/assembler/packetconn"

	_ "github.com/gzjjjfree/v5-result/transport/internet/request/stereotype/meek"
	_ "github.com/gzjjjfree/v5-result/transport/internet/request/stereotype/mekya"

	_ "github.com/gzjjjfree/v5-result/transport/internet/request/roundtripperreverserserver/clicommand"

	_ "github.com/gzjjjfree/v5-result/transport/internet/dtls"

	_ "github.com/gzjjjfree/v5-result/transport/internet/httpupgrade"

	_ "github.com/gzjjjfree/v5-result/transport/internet/hysteria2"

	_ "github.com/gzjjjfree/v5-result/transport/internet/tlsmirror/mirrorenrollment/roundtripperenrollmentconfirmation"
	_ "github.com/gzjjjfree/v5-result/transport/internet/tlsmirror/server"

	// Transport headers
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/http"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/noop"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/srtp"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/tls"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/utp"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/wechat"
	_ "github.com/gzjjjfree/v5-result/transport/internet/headers/wireguard"

	// Geo loaders
	_ "github.com/gzjjjfree/v5-result/infra/conf/geodata/memconservative"
	_ "github.com/gzjjjfree/v5-result/infra/conf/geodata/standard"

	// JSON, TOML, YAML config support. (jsonv4) This disable selective compile
	_ "github.com/gzjjjfree/v5-result/main/formats"

	// commands
	_ "github.com/gzjjjfree/v5-result/main/commands/all"

	// engineering commands
	_ "github.com/gzjjjfree/v5-result/main/commands/all/engineering"
	_ "github.com/gzjjjfree/v5-result/main/commands/all/engineering/generateRandomData"

	// Commands that rely on jsonv4 format This disable selective compile
	_ "github.com/gzjjjfree/v5-result/main/commands/all/api/jsonv4"
	_ "github.com/gzjjjfree/v5-result/main/commands/all/jsonv4"

	// V5 version of json configure file parser
	_ "github.com/gzjjjfree/v5-result/infra/conf/v5cfg"

	// Simplified config
	_ "github.com/gzjjjfree/v5-result/proxy/http/simplified"
	_ "github.com/gzjjjfree/v5-result/proxy/shadowsocks/simplified"
	_ "github.com/gzjjjfree/v5-result/proxy/socks/simplified"
	_ "github.com/gzjjjfree/v5-result/proxy/trojan/simplified"

	// Subscription Supports
	_ "github.com/gzjjjfree/v5-result/app/subscription/subscriptionmanager"

	_ "github.com/gzjjjfree/v5-result/app/subscription/subscriptionmanager/command"

	// Subscription Containers: general purpose
	_ "github.com/gzjjjfree/v5-result/app/subscription/containers/base64urlline"
	_ "github.com/gzjjjfree/v5-result/app/subscription/containers/dataurlsingle"
	_ "github.com/gzjjjfree/v5-result/app/subscription/containers/jsonfieldarray"
	_ "github.com/gzjjjfree/v5-result/app/subscription/containers/jsonfieldarray/jsonified"
	_ "github.com/gzjjjfree/v5-result/app/subscription/containers/urlline"

	// Subscription Fetchers
	_ "github.com/gzjjjfree/v5-result/app/subscription/documentfetcher/dataurlfetcher"
	_ "github.com/gzjjjfree/v5-result/app/subscription/documentfetcher/httpfetcher"

	// Subscription Entries Converters
	_ "github.com/gzjjjfree/v5-result/app/subscription/entries/nonnative"
	_ "github.com/gzjjjfree/v5-result/app/subscription/entries/outbound" // Natively Supported Outbound Format
)
