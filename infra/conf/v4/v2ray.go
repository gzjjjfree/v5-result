package v4

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	"github.com/v2fly/v2ray-core/v5/app/stats"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/features"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon/loader"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon/muxcfg"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon/proxycfg"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon/sniffer"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/dns"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/log"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/router"
	"github.com/v2fly/v2ray-core/v5/infra/conf/v5cfg"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

var (
	inboundConfigLoader = loader.NewJSONConfigLoader(loader.ConfigCreatorCache{
		"dokodemo-door": func() interface{} { return new(DokodemoConfig) },
		"http":          func() interface{} { return new(HTTPServerConfig) },
		"shadowsocks":   func() interface{} { return new(ShadowsocksServerConfig) },
		"socks":         func() interface{} { return new(SocksServerConfig) },
		"vless":         func() interface{} { return new(VLessInboundConfig) },
		"vmess":         func() interface{} { return new(VMessInboundConfig) },
		"trojan":        func() interface{} { return new(TrojanServerConfig) },
		"hysteria2":     func() interface{} { return new(Hysteria2ServerConfig) },
	}, "protocol", "settings")

	outboundConfigLoader = loader.NewJSONConfigLoader(loader.ConfigCreatorCache{
		"blackhole":   func() interface{} { return new(BlackholeConfig) },
		"freedom":     func() interface{} { return new(FreedomConfig) },
		"http":        func() interface{} { return new(HTTPClientConfig) },
		"shadowsocks": func() interface{} { return new(ShadowsocksClientConfig) },
		"socks":       func() interface{} { return new(SocksClientConfig) },
		"vless":       func() interface{} { return new(VLessOutboundConfig) },
		"vmess":       func() interface{} { return new(VMessOutboundConfig) },
		"trojan":      func() interface{} { return new(TrojanClientConfig) },
		"hysteria2":   func() interface{} { return new(Hysteria2ClientConfig) },
		"dns":         func() interface{} { return new(DNSOutboundConfig) },
		"loopback":    func() interface{} { return new(LoopbackConfig) },
	}, "protocol", "settings")
)

func toProtocolList(s []string) ([]proxyman.KnownProtocols, error) {
	kp := make([]proxyman.KnownProtocols, 0, 8)
	for _, p := range s {
		switch strings.ToLower(p) {
		case "http":
			kp = append(kp, proxyman.KnownProtocols_HTTP)
		case "https", "tls", "ssl":
			kp = append(kp, proxyman.KnownProtocols_TLS)
		default:
			return nil, newError("Unknown protocol: ", p)
		}
	}
	return kp, nil
}

type InboundDetourAllocationConfig struct {
	Strategy    string  `json:"strategy"`
	Concurrency *uint32 `json:"concurrency"`
	RefreshMin  *uint32 `json:"refresh"`
}

// Build implements Buildable.
func (c *InboundDetourAllocationConfig) Build() (*proxyman.AllocationStrategy, error) {
	config := new(proxyman.AllocationStrategy)
	switch strings.ToLower(c.Strategy) {
	case "always":
		config.Type = proxyman.AllocationStrategy_Always
	case "random":
		config.Type = proxyman.AllocationStrategy_Random
	case "external":
		config.Type = proxyman.AllocationStrategy_External
	default:
		return nil, newError("unknown allocation strategy: ", c.Strategy)
	}
	if c.Concurrency != nil {
		config.Concurrency = &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
			Value: *c.Concurrency,
		}
	}

	if c.RefreshMin != nil {
		config.Refresh = &proxyman.AllocationStrategy_AllocationStrategyRefresh{
			Value: *c.RefreshMin,
		}
	}

	return config, nil
}

type InboundDetourConfig struct {
	Protocol       string                         `json:"protocol"`
	PortRange      *cfgcommon.PortRange           `json:"port"`
	ListenOn       *cfgcommon.Address             `json:"listen"`
	Settings       *json.RawMessage               `json:"settings"`
	Tag            string                         `json:"tag"`
	Allocation     *InboundDetourAllocationConfig `json:"allocate"`
	StreamSetting  *StreamConfig                  `json:"streamSettings"`
	DomainOverride *cfgcommon.StringList          `json:"domainOverride"`
	SniffingConfig *sniffer.SniffingConfig        `json:"sniffing"`
}

// Build implements Buildable.
func (c *InboundDetourConfig) Build() (*core.InboundHandlerConfig, error) {
	receiverSettings := &proxyman.ReceiverConfig{}

	if c.ListenOn == nil {
		// Listen on anyip, must set PortRange
		if c.PortRange == nil {
			return nil, newError("Listen on AnyIP but no Port(s) set in InboundDetour.")
		}
		receiverSettings.PortRange = c.PortRange.Build()
	} else {
		// Listen on specific IP or Unix Domain Socket
		receiverSettings.Listen = c.ListenOn.Build()
		listenDS := c.ListenOn.Family().IsDomain() && (filepath.IsAbs(c.ListenOn.Domain()) || c.ListenOn.Domain()[0] == '@')
		listenIP := c.ListenOn.Family().IsIP() || (c.ListenOn.Family().IsDomain() && c.ListenOn.Domain() == "localhost")
		switch {
		case listenIP:
			// Listen on specific IP, must set PortRange
			if c.PortRange == nil {
				return nil, newError("Listen on specific ip without port in InboundDetour.")
			}
			// Listen on IP:Port
			receiverSettings.PortRange = c.PortRange.Build()
		case listenDS:
			if c.PortRange != nil {
				// Listen on Unix Domain Socket, PortRange should be nil
				receiverSettings.PortRange = nil
			}
		default:
			return nil, newError("unable to listen on domain address: ", c.ListenOn.Domain())
		}
	}

	//c.Allocation = nil // 屏蔽
	if c.Allocation != nil {
		concurrency := -1
		if c.Allocation.Concurrency != nil && c.Allocation.Strategy == "random" {
			concurrency = int(*c.Allocation.Concurrency)
		}
		portRange := int(c.PortRange.To - c.PortRange.From + 1)
		if concurrency >= 0 && concurrency >= portRange {
			return nil, newError("not enough ports. concurrency = ", concurrency, " ports: ", c.PortRange.From, " - ", c.PortRange.To)
		}

		as, err := c.Allocation.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.AllocationStrategy = as
	}

	//c.StreamSetting = nil // 屏蔽
	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.StreamSettings = ss
	}
	if c.SniffingConfig != nil {
		s, err := c.SniffingConfig.Build()
		if err != nil {
			return nil, newError("failed to build sniffing config").Base(err)
		}
		receiverSettings.SniffingSettings = s
	}

	//c.DomainOverride = nil // 屏蔽
	if c.DomainOverride != nil {
		kp, err := toProtocolList(*c.DomainOverride)
		if err != nil {
			return nil, newError("failed to parse inbound detour config").Base(err)
		}
		receiverSettings.DomainOverride = kp
	}

	//c.Settings = nil // 屏蔽
	settings := []byte("{}")
	if c.Settings != nil {
		settings = ([]byte)(*c.Settings)
	}
	rawConfig, err := inboundConfigLoader.LoadWithID(settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound detour config.").Base(err)
	}
	if dokodemoConfig, ok := rawConfig.(*DokodemoConfig); ok {
		receiverSettings.ReceiveOriginalDestination = dokodemoConfig.Redirect
	}
	ts, err := rawConfig.(cfgcommon.Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.InboundHandlerConfig{
		Tag:              c.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverSettings),
		ProxySettings:    serial.ToTypedMessage(ts),
	}, nil
}

type OutboundDetourConfig struct {
	Protocol       string                `json:"protocol"`
	SendThrough    *cfgcommon.Address    `json:"sendThrough"`
	Tag            string                `json:"tag"`
	Settings       *json.RawMessage      `json:"settings"`
	StreamSetting  *StreamConfig         `json:"streamSettings"`
	ProxySettings  *proxycfg.ProxyConfig `json:"proxySettings"`
	MuxSettings    *muxcfg.MuxConfig     `json:"mux"`
	DomainStrategy string                `json:"domainStrategy"`
}

// Build implements Buildable.
func (c *OutboundDetourConfig) Build() (*core.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if c.SendThrough != nil {
		address := c.SendThrough
		if address.Family().IsDomain() {
			return nil, newError("unable to send through: " + address.String())
		}
		senderSettings.Via = address.Build()
	}

	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}

	if c.ProxySettings != nil {
		ps, err := c.ProxySettings.Build()
		if err != nil {
			return nil, newError("invalid outbound detour proxy settings.").Base(err)
		}
		senderSettings.ProxySettings = ps
	}

	if c.MuxSettings != nil {
		senderSettings.MultiplexSettings = c.MuxSettings.Build()
	}

	senderSettings.DomainStrategy = proxyman.SenderConfig_AS_IS
	switch strings.ToLower(c.DomainStrategy) {
	case "useip", "use_ip", "use-ip":
		senderSettings.DomainStrategy = proxyman.SenderConfig_USE_IP
	case "useip4", "useipv4", "use_ip4", "use_ipv4", "use_ip_v4", "use-ip4", "use-ipv4", "use-ip-v4":
		senderSettings.DomainStrategy = proxyman.SenderConfig_USE_IP4
	case "useip6", "useipv6", "use_ip6", "use_ipv6", "use_ip_v6", "use-ip6", "use-ipv6", "use-ip-v6":
		senderSettings.DomainStrategy = proxyman.SenderConfig_USE_IP6
	}

	settings := []byte("{}")
	if c.Settings != nil {
		settings = ([]byte)(*c.Settings)
	}
	rawConfig, err := outboundConfigLoader.LoadWithID(settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to parse to outbound detour config.").Base(err)
	}

	ts, err := rawConfig.(cfgcommon.Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		Tag:            c.Tag,
		ProxySettings:  serial.ToTypedMessage(ts),
	}, nil
}

type StatsConfig struct{}

// Build implements Buildable.
func (c *StatsConfig) Build() (*stats.Config, error) {
	return &stats.Config{}, nil
}

type Config struct {
	// Port of this Point server.
	// Deprecated: Port exists for historical compatibility
	// and should not be used.
	Port uint16 `json:"port"`

	// Deprecated: InboundConfig exists for historical compatibility
	// and should not be used.
	InboundConfig *InboundDetourConfig `json:"inbound"`

	// Deprecated: OutboundConfig exists for historical compatibility
	// and should not be used.
	OutboundConfig *OutboundDetourConfig `json:"outbound"`

	// Deprecated: InboundDetours exists for historical compatibility
	// and should not be used.
	InboundDetours []InboundDetourConfig `json:"inboundDetour"`

	// Deprecated: OutboundDetours exists for historical compatibility
	// and should not be used.
	OutboundDetours []OutboundDetourConfig `json:"outboundDetour"`

	LogConfig        *log.LogConfig          `json:"log"`
	RouterConfig     *router.RouterConfig    `json:"routing"`
	DNSConfig        *dns.DNSConfig          `json:"dns"`
	InboundConfigs   []InboundDetourConfig   `json:"inbounds"`
	OutboundConfigs  []OutboundDetourConfig  `json:"outbounds"`
	Transport        *TransportConfig        `json:"transport"`
	Policy           *PolicyConfig           `json:"policy"`
	API              *APIConfig              `json:"api"`
	Stats            *StatsConfig            `json:"stats"`
	Reverse          *ReverseConfig          `json:"reverse"`
	FakeDNS          *dns.FakeDNSConfig      `json:"fakeDns"`
	BrowserForwarder *BrowserForwarderConfig `json:"browserForwarder"`
	Observatory      *ObservatoryConfig      `json:"observatory"`
	BurstObservatory *BurstObservatoryConfig `json:"burstObservatory"`
	MultiObservatory *MultiObservatoryConfig `json:"multiObservatory"`

	Services map[string]*json.RawMessage `json:"services"`
}

func (c *Config) findInboundTag(tag string) int {
	found := -1
	for idx, ib := range c.InboundConfigs {
		if ib.Tag == tag {
			found = idx
			break
		}
	}
	return found
}

func (c *Config) findOutboundTag(tag string) int {
	found := -1
	for idx, ob := range c.OutboundConfigs {
		if ob.Tag == tag {
			found = idx
			break
		}
	}
	return found
}

func applyTransportConfig(s *StreamConfig, t *TransportConfig) {
	if s.TCPSettings == nil {
		s.TCPSettings = t.TCPConfig
	}
	if s.KCPSettings == nil {
		s.KCPSettings = t.KCPConfig
	}
	if s.WSSettings == nil {
		s.WSSettings = t.WSConfig
	}
	if s.HTTPSettings == nil {
		s.HTTPSettings = t.HTTPConfig
	}
	if s.DSSettings == nil {
		s.DSSettings = t.DSConfig
	}
}

// Build implements Buildable.
func (c *Config) Build() (*core.Config, error) {
	if err := PostProcessConfigureFile(c); err != nil {
		return nil, err
	}

	config := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	//c.API = nil // 屏蔽
	if c.API != nil {
		apiConf, err := c.API.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(apiConf))
	}

	//c.Stats = nil // 屏蔽
	if c.Stats != nil {
		statsConf, err := c.Stats.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(statsConf))
	}

	//c.LogConfig = nil // 屏蔽
	var logConfMsg *anypb.Any
	if c.LogConfig != nil {
		logConfMsg = serial.ToTypedMessage(c.LogConfig.Build())
	} else {
		logConfMsg = serial.ToTypedMessage(log.DefaultLogConfig())
	}
	// let logger module be the first App to start,
	// so that other modules could print log during initiating
	config.App = append([]*anypb.Any{logConfMsg}, config.App...)

	if c.RouterConfig != nil {
		routerConfig, err := c.RouterConfig.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(routerConfig))
	}

	if c.FakeDNS != nil {
		features.PrintDeprecatedFeatureWarning("root fakedns settings")
		if c.DNSConfig != nil {
			c.DNSConfig.FakeDNS = c.FakeDNS
		} else {
			c.DNSConfig = &dns.DNSConfig{
				FakeDNS: c.FakeDNS,
			}
		}
	}

	if c.DNSConfig != nil {
		dnsApp, err := c.DNSConfig.Build()
		if err != nil {
			return nil, newError("failed to parse DNS config").Base(err)
		}
		config.App = append(config.App, serial.ToTypedMessage(dnsApp))
	}

	if c.Policy != nil {
		pc, err := c.Policy.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(pc))
	}

	if c.Reverse != nil {
		r, err := c.Reverse.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(r))
	}

	if c.BrowserForwarder != nil {
		r, err := c.BrowserForwarder.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(r))
	}

	if c.Observatory != nil {
		r, err := c.Observatory.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(r))
	}

	if c.BurstObservatory != nil {
		r, err := c.BurstObservatory.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(r))
	}

	if c.MultiObservatory != nil {
		r, err := c.MultiObservatory.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(r))
	}

	// Load Additional Services that do not have a json translator

	for serviceName, service := range c.Services {
		servicePackedConfig, err := v5cfg.LoadHeterogeneousConfigFromRawJSON(context.Background(), "service", serviceName, *service)
		if err != nil {
			return nil, newError(fmt.Sprintf("failed to parse %v config in Services", serviceName)).Base(err)
		}
		config.App = append(config.App, serial.ToTypedMessage(servicePackedConfig))
	}

	var inbounds []InboundDetourConfig

	if c.InboundConfig != nil {
		inbounds = append(inbounds, *c.InboundConfig)
	}

	if len(c.InboundDetours) > 0 {
		inbounds = append(inbounds, c.InboundDetours...)
	}

	if len(c.InboundConfigs) > 0 {
		inbounds = append(inbounds, c.InboundConfigs...)
	}

	// Backward compatibility.
	if len(inbounds) > 0 && inbounds[0].PortRange == nil && c.Port > 0 {
		inbounds[0].PortRange = &cfgcommon.PortRange{
			From: uint32(c.Port),
			To:   uint32(c.Port),
		}
	}

	for _, rawInboundConfig := range inbounds {
		if c.Transport != nil {
			if rawInboundConfig.StreamSetting == nil {
				rawInboundConfig.StreamSetting = &StreamConfig{}
			}
			applyTransportConfig(rawInboundConfig.StreamSetting, c.Transport)
		}
		ic, err := rawInboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic)
	}

	var outbounds []OutboundDetourConfig

	if c.OutboundConfig != nil {
		outbounds = append(outbounds, *c.OutboundConfig)
	}

	if len(c.OutboundDetours) > 0 {
		outbounds = append(outbounds, c.OutboundDetours...)
	}

	if len(c.OutboundConfigs) > 0 {
		outbounds = append(outbounds, c.OutboundConfigs...)
	}

	for _, rawOutboundConfig := range outbounds {
		if c.Transport != nil {
			if rawOutboundConfig.StreamSetting == nil {
				rawOutboundConfig.StreamSetting = &StreamConfig{}
			}
			applyTransportConfig(rawOutboundConfig.StreamSetting, c.Transport)
		}		

		// 逐级检查指针，防止出现 nil pointer dereference 导致崩溃
		if rawOutboundConfig.StreamSetting != nil &&
			rawOutboundConfig.StreamSetting.TLSSettings != nil {

			// 检查是否非空且长度大于 0
			if len(rawOutboundConfig.StreamSetting.TLSSettings.ECHDOHServer) > 0 &&
				rawOutboundConfig.StreamSetting.TLSSettings.ServerName != "" {

				fmt.Println("rawOutboundConfig.Tag:", rawOutboundConfig.Tag)
				fmt.Println("[ECH] DOH 服务器配置:", string(rawOutboundConfig.StreamSetting.TLSSettings.ECHDOHServer))
				fmt.Println("[ECH] ServerName 配置:", rawOutboundConfig.StreamSetting.TLSSettings.ServerName)

				addr := net.ParseAddress(rawOutboundConfig.StreamSetting.TLSSettings.ServerName)
				if addr.Family().IsDomain() {
					fmt.Println("[ECH] ServerName 是域名，正在解析...")
					go tls.GetECHConfigBackground(addr.String(), rawOutboundConfig.StreamSetting.TLSSettings.ECHDOHServer)
				}
			} else {
				fmt.Println("[ECH] 配置不存在或为空字节数组")
			}
		}

		// 出站 tag 以 "cdn-" 开头时，以 IP 池地址建立出站列表
		if strings.HasPrefix(strings.ToLower(rawOutboundConfig.Tag), "cdn-") {
			d, derr := Configloads()
			if derr == nil && len(d) > 0 {
				if len(d) > 50 {
					d = d[:50]
				}
				// 保存原始 Tag 模板，防止累加导致的 Tag 错误
				originalTag := rawOutboundConfig.Tag

				for i := 0; i < len(d); i++ {
					// 拷贝配置，防止修改影响后续循环
					tempConfig := rawOutboundConfig

					// 生成规范的唯一 Tag，例如 cdn-node-0, cdn-node-1
					tempConfig.Tag = fmt.Sprintf("%v-%d", originalTag, i)

					// 判断协议类型并修改对应的 Address
					settings := []byte("{}")
					if tempConfig.Settings != nil {
						settings = ([]byte)(*tempConfig.Settings)
					}

					// 我们直接操作 JSON Map，绕过结构体的 Marshal 限制
					var settingsMap map[string]any
					json.Unmarshal(settings, &settingsMap)

					switch tempConfig.Protocol {
					case "vless":
						// 如果是 VLess，修改第一个 vnext 的地址
						if vnext, ok := settingsMap["vnext"].([]any); ok && len(vnext) > 0 {
							if firstVnext, ok := vnext[0].(map[string]any); ok {
								firstVnext["address"] = d[i].Addresses.String()
							}
						}
					case "vmess":
						// 如果是 VMess，修改第一个 receiver 的地址
						if vmess, ok := settingsMap["Receivers"].([]any); ok && len(vmess) > 0 {
							if firstvmess, ok := vmess[0].(map[string]interface{}); ok {
								firstvmess["address"] = d[i].Addresses.String()
							}
						}
					case "trojan":
						// 如果是 Trojan 或 Shadowsocks 同字段 Servers
						if trojan, ok := settingsMap["Servers"].([]any); ok && len(trojan) > 0 {
							if firsttrojan, ok := trojan[0].(map[string]any); ok {
								firsttrojan["Servers"] = d[i].Addresses.String()
							}
						}
					default:
						// 如果是其他协议，可以尝试通过反射或忽略
						fmt.Printf("警告: 协议 %s 暂不支持动态注入优选 IP\n", tempConfig.Protocol)
					}

					// 将修改后的 Map 封回 Settings
					modifiedSettings, err := json.Marshal(settingsMap)
					if err != nil {
						return nil, err
					}
					tempConfig.Settings = (*json.RawMessage)(&modifiedSettings)

					oc, err := tempConfig.Build()
					if err != nil {
						return nil, err
					}
					config.Outbound = append(config.Outbound, oc)
				}
				// 成功处理完优选列表后，跳过原本模板的构建（continue）
				continue
			}
		}

		oc, err := rawOutboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Outbound = append(config.Outbound, oc)
	}

	return config, nil
}

type vaddresses struct {
	Addresses *cfgcommon.Address `json:"address"`
}

func Configloads() ([]vaddresses, error) {
	// 获取二进制文件所在的绝对路径，而不是依赖执行时的终端位置
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(exePath)
	dirPath := filepath.Join(baseDir, "result") // 这样无论在哪运行，都会找二进制旁边的 result 文件夹

	var files []string

	// 检查文件夹是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("[ERROR] 文件夹不存在: %s\n", dirPath)
		dirPath = filepath.Join(".", "result")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("[ERROR] 文件夹不存在: %s\n", dirPath)
			return nil, err
		}
		fmt.Printf("读取文件夹: %s\n", dirPath)
	}
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("WalkDir err: %v", err.Error())
			return err // 如果遍历过程中出现错误，则返回错误
		}
		if !d.IsDir() && strings.Contains(d.Name(), "result") { // 检查是否为文件且文件名包含 "result"
			files = append(files, path)
			fmt.Printf("files: %v\n", files)
		}
		return nil
	})

	var addresses []vaddresses
	for _, file := range files {
		content, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue // 继续处理下一个文件
		}
		defer content.Close()

		var addr []vaddresses
		decoder := json.NewDecoder(content)
		err = decoder.Decode(&addr)
		if err != nil {
			fmt.Println("解析 result.json 文件出错: ", err)
			return nil, err
		}
		addresses = append(addresses, addr...)
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(addresses), func(i, j int) {
		addresses[i], addresses[j] = addresses[j], addresses[i]
	})
	return addresses, nil
}
