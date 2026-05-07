//go:build go1.23
// +build go1.23

package tls

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"

	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
)

var (
	// 保护 statusMap 的并发安全
	statusMutex   sync.Mutex
	applyEchMutex sync.Mutex
	// 存储每个域名的状态：如果正在更新，对应 chan 会有数据或被关闭
	updatingState = make(map[string]chan struct{})

	UpdateSignal = make(chan string, 1)
	// Map：key 是 domain，value 是对应的 ECH 配置字节
	globalEchCache = make(map[string][]byte)
)

// GetECHConfigBackground 后台监听函数
// sig: 信号通道
func GetECHConfigBackground(serverName string, workerDomain string, addresses []string) {
	fmt.Println("[ECH] 后台同步协程已启动，等待信号...")
	
	// 启动时先主动同步一次，确保初始状态是最新的
	//doUpdate(ctx, serverName, workerDomain, addresses)

	// 封装一个带锁的更新函数
	safeUpdate := func(target string) {
		// 为每次更新创建独立的 10 秒超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		statusMutex.Lock()
		ch := make(chan struct{})
		updatingState[target] = ch
		statusMutex.Unlock()

		// 使用 defer 确保无论 doUpdate 是否报错，都会解除阻塞
		defer func() {
			statusMutex.Lock()
			delete(updatingState, target)
			close(ch)
			statusMutex.Unlock()
		}()

		doUpdate(ctx, target, workerDomain, addresses)
	}

	// 启动初始同步
	safeUpdate(serverName)

	// 监听信号，UpdateSignal 此时传递的是需要更新的 domain
	for domain := range UpdateSignal {
		// 如果信号传递的域名匹配或有通用更新逻辑
		safeUpdate(domain)
		fmt.Printf("[ECH] 域名 %s 任务完成，继续等待...\n", domain)
	}
}

// 执行更新操作
func doUpdate(ctx context.Context, serverName string, workerDomain string, addresses []string) {
	fmt.Println("[ECH] 收到更新信号，正在拉取最新配置...")
	// 注意：这里的 workerDomain 是你分配给这个 Worker 的域名（必须是自定义域名）

	// 定义备选 IP 列表
	ips := []string{
		"104.16.123.99:443",
	}

	// 遍历并补全端口，然后追加到 ips 列表中
	for _, addr := range addresses {
		// 如果地址里没写端口，手动补上 :443
		if !strings.Contains(addr, ":") {
			addr = addr + ":443"
		}
		ips = append(ips, addr)
	}

	targetDomain := serverName

	// 构造完整的请求 URL
	// curl.exe -v --resolve 自定义域名:443:104.16.123.99 "https://自定义域名/getech?domain=serverName"
	finalUrl := fmt.Sprintf("https://%s/getech?domain=%s", workerDomain, targetDomain)

	// 2. 开始轮询重试
	var body []byte
	var success bool

	for _, cfIP := range ips {
		fmt.Printf("[ECH] 尝试使用 IP: %s ...拉取域名: %s\n", cfIP, targetDomain)

		transport := &http.Transport{
			DialContext: func(dialCtx context.Context, network, addr string) (net.Conn, error) {
				// 设置拨号超时，防止在单个 IP 上卡死太久
				dialer := &net.Dialer{Timeout: 3 * time.Second}
				return dialer.DialContext(ctx, network, cfIP)
			},
			TLSClientConfig: &tls.Config{
				ServerName: workerDomain,
				MinVersion: tls.VersionTLS13,
			},
			// 每次使用新 IP 都要禁用长连接重用，确保拨号发生
			DisableKeepAlives: true,
		}

		client := &http.Client{Transport: transport, Timeout: 5 * time.Second}

		resp, err := client.Get(finalUrl)
		if err != nil {
			fmt.Printf("[ECH] IP %s 拉取域名: %s 连接失败: %v\n", cfIP, targetDomain, err)
			continue // 尝试下一个 IP
		}

		if resp.StatusCode == http.StatusOK {
			var readErr error
			body, readErr = io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr == nil {
				success = true
				fmt.Printf("[ECH] IP %s 拉取成功！拉取域名: %s\n", cfIP, targetDomain)
				break // 成功获取，跳出循环
			}
		} else {
			fmt.Printf("[ECH] IP %s 拉取域名: %s 返回状态码: %d\n", cfIP, targetDomain, resp.StatusCode)
			resp.Body.Close()
		}
	}

	if !success {
		fmt.Printf("[ECH] 所有备选 IP 均已尝试，拉取全部失败拉取域名: %s\n", targetDomain)
		return
	}

	content := strings.TrimSpace(string(body))
	if content == "" {
		fmt.Printf("[ECH] 获取内容为空拉取域名: %s\n", targetDomain)
		return
	}

	ECHConfigBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		fmt.Printf("[ECH] Base64 解码失败: %v, 内容: [%s]\n", err, content)
		return
	}

	// 更新内存
	applyEchMutex.Lock()
	// 使用 serverName 作为 Key 存储
	globalEchCache[serverName] = ECHConfigBytes
	applyEchMutex.Unlock()

	fmt.Printf("[ECH] 内存配置已更新，长度: %d 拉取域名: %s\n", len(ECHConfigBytes), serverName)
}

func ApplyECH(c *Config, config *tls.Config) error {
	// 检查该域名是否正在更新
	statusMutex.Lock()
	waitCh, isUpdating := updatingState[c.ServerName]
	statusMutex.Unlock()

	if isUpdating {
		fmt.Printf("[ECH] 域名 %s 正在更新，请求进入阻塞等待...\n", c.ServerName)
		<-waitCh // 阻塞在这里，直到 doUpdate 完成并 close(ch)
	}

	// 此时更新已完成，安全读取数据
	// 注意：这里的读取最好也加一下你原来的全局轻量锁，或者确保 map 写入后不再修改
	applyEchMutex.Lock()
	config.EncryptedClientHelloConfigList = globalEchCache[c.ServerName]
	applyEchMutex.Unlock()

	return nil
}

//func ApplyECH(c *Config, config *tls.Config) error { //
//	var ECHConfig []byte
//	var err error
//	var domain string
//
//	if len(c.EchConfig) > 0 {
//		ECHConfig = c.EchConfig
//	} else { // ECH config > DOH lookup
//		if c.EchQueryDomain == "" {
//			domain = config.ServerName
//		} else {
//			domain = c.EchQueryDomain
//		}
//		addr := net.ParseAddress(domain)
//		if !addr.Family().IsDomain() {
//			return newError("Using DOH for ECH needs SNI")
//		}
//		ECHConfig, err = QueryRecord(addr.Domain(), c.Ech_DOHserver)
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//	}
//
//	config.EncryptedClientHelloConfigList = ECHConfig
//	return nil
//}

type record struct {
	record []byte
	expire time.Time
}

var (
	dnsCache = make(map[string]record)
	mutex    sync.RWMutex
)

func QueryRecord(domain string, server string) ([]byte, error) {
	mutex.Lock()
	rec, found := dnsCache[domain]
	if found && rec.expire.After(time.Now()) {
		mutex.Unlock()
		return rec.record, nil
	}
	mutex.Unlock()

	newError("Trying to query ECH config for domain: ", domain, " with ECH server: ", server).AtDebug().WriteToLog()
	record, ttl, err := dohQuery(server, domain)
	if err != nil {
		return []byte{}, err
	}

	if ttl < 600 {
		ttl = 600
	}

	mutex.Lock()
	defer mutex.Unlock()
	rec.record = record
	rec.expire = time.Now().Add(time.Second * time.Duration(ttl))
	dnsCache[domain] = rec
	return record, nil
}

func dohQuery(server string, domain string) ([]byte, uint32, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeHTTPS)
	m.Id = 0
	msg, err := m.Pack()
	if err != nil {
		return []byte{}, 0, err
	}
	tr := &http.Transport{
		IdleConnTimeout:   90 * time.Second,
		ForceAttemptHTTP2: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dest, err := net.ParseDestination(network + ":" + addr)
			if err != nil {
				return nil, err
			}
			conn, err := internet.DialSystem(ctx, dest, nil)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	req, err := http.NewRequest("POST", server, bytes.NewReader(msg))
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/dns-message")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, 0, newError("query failed with response code:", resp.StatusCode)
	}
	respMsg := new(dns.Msg)
	err = respMsg.Unpack(respBody)
	if err != nil {
		return []byte{}, 0, err
	}
	if len(respMsg.Answer) > 0 {
		for _, answer := range respMsg.Answer {
			if https, ok := answer.(*dns.HTTPS); ok && https.Hdr.Name == dns.Fqdn(domain) {
				for _, v := range https.Value {
					if echConfig, ok := v.(*dns.SVCBECHConfig); ok {
						newError(context.Background(), "Get ECH config:", echConfig.String(), " TTL:", respMsg.Answer[0].Header().Ttl).AtDebug().WriteToLog()
						return echConfig.ECH, answer.Header().Ttl, nil
					}
				}
			}
		}
	}
	return []byte{}, 0, newError("no ech record found")
}
