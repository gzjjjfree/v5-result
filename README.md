<div>
  <img width="190" height="210" align="left" src="https://raw.githubusercontent.com/v2fly/v2fly-github-io/master/docs/.vuepress/public/readme-logo.png" alt="V2Ray"/>
  <br>
  <h1>Project V</h1>
  <p>Project V is a set of network tools that helps you to build your own computer network. It secures your network connections and thus protects your privacy.</p>
</div>

![Downloads](https://img.shields.io/github/downloads/gzjjjfree/v5-result/total?style=flat-square)

# V2Ray v5-result 定制版

本分支基于 V2Ray v5 开发，针对 Cloudflare 优选 IP 场景进行了深度定制。核心特性包括：**自动加载优选 IP 池**、**动态生成负载均衡出站**以及**基于优选 IP 的 ECH (Encrypted Client Hello) 自动化更新机制**。

---

## 核心特性

1.  **自动 IP 池加载**：自动读取同级目录 `result/result.json` 中的 IP 地址。
2.  **动态出站生成**：只要出站 `tag` 以 `cdn-` 开头（如 `cdn-vless`），程序会自动克隆该配置，并为 IP 池中的每个 IP 生成独立的出站节点（如 `cdn-vless-0`, `cdn-vless-1` 等）。
3.  **ECH 自动化热更新**：
    * 通过优选 IP 列表轮询访问自定义 Workers。
    * 自动获取并更新内存中的 ECH 配置，无需重启。
    * 支持 `ServerName` 级别的 ECH 匹配。

---

## 文件目录结构

请确保 `result.json` 放置在如下位置：

```text
.
├── v2ray (可执行文件)
├── config.json
└── result/
    └── result.json
```

---

## 配置文件指南

### 1. 准备 IP 池 (`result/result.json`)
请将优选工具扫描出的结果按以下格式存入：

```json
[
    {
        "address": "104.16.244.51"
    },
    {
        "address": "104.19.46.94"
    }
]
```

### 2. 配置主文件 (`config.json`)

#### 出站配置 (Outbounds)
只需配置一个模板，`tag` 必须以 `cdn-` 为前缀。

```json
"outbounds": [
    {
        "protocol": "vless", 
        "tag": "cdn-vless",
        "settings": { 
            "vnext": [
                {
                    "port": 443, 
                    "users": [
                        {
                            "id": "你的UUID",
                            "encryption": "none"
                        }
                    ]
                }
            ]
        },
        "streamSettings": { 
            "network": "ws",
            "security": "tls", 
            "tlsSettings": {
                "serverName": "cf 代理你的网站名",
                "echDohServer": "cf 请求 ECH 的 worker 的自定义域名，不要带 https://", 
                "allowInsecure": false
            },
            "wsSettings": {
                "path": "/",
                "headers": {
                    "Host": "cf 代理你的网站名"
                }
            }
        },
        "mux": { "enabled": false }
    }
]
```

#### 请求 ECH 的 workers 代码

```workers
export default {
  async fetch(request) {
    const url = new URL(request.url);
    const domain = url.searchParams.get('domain');

    if (domain) {
      try {
        const dohUrl = `https://1.1.1.1/dns-query?name=${domain}&type=65`;
        const response = await fetch(dohUrl, {
          headers: { "accept": "application/dns-json" }
        });
        const json = await response.json();
        
        const echConfig = extractEch(json);
        if (echConfig) {
          return new Response(echConfig, { headers: { "Content-Type": "text/plain" } });
        }
        return new Response("ECH not found in record", { status: 404 });
      } catch (e) {
        return new Response("Error: " + e.message, { status: 500 });
      }
    }
    return new Response("Missing domain", { status: 400 });
  }
};

function extractEch(dnsJson) {
  if (!dnsJson.Answer) return null;

  for (const record of dnsJson.Answer) {
    if (record.type === 65) {
      const data = record.data;

      // 情况 A: 已经是易读格式 ech="xxx"
      const match = data.match(/ech="([^"]+)"/);
      if (match) return match[1];

      // 情况 B: 十六进制格式 \# <len> <hex_data>
      if (data.startsWith("\\#")) {
        // 移除前缀 "\# 136 "（具体的长度数字可能不同）
        const parts = data.split(' ');
        // 真正的十六进制数据从索引 2 或 3 开始，我们将所有部分合并
        const hex = parts.slice(2).join('');
        
        // ECH 的标识符是 0005
        const echIndex = hex.indexOf("0005");
        if (echIndex !== -1) {
          // 0005 后面是 2 字节长度 (4个字符)
          const lenHex = hex.substring(echIndex + 4, echIndex + 8);
          const len = parseInt(lenHex, 16);
          // 提取 ECH 核心数据
          const echHex = hex.substring(echIndex + 8, echIndex + 8 + (len * 2));
          
          // 将 Hex 转换为 Uint8Array，再转为 Base64
          const bytes = new Uint8Array(echHex.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));
          return btoa(String.fromCharCode(...bytes));
        }
      }
    }
  }
  return null;
}
```

#### 路由与负载均衡 (Routing)
使用 `balancer` 自动匹配所有动态生成的 `cdn-` 节点。

```json
"routing": {
    "domainStrategy": "AsIs",  
    "balancers": [
        {
            "tag": "cdn-balancer",
            "selector": ["cdn-"],
            "strategy": {
                "type": "random" 
            }
        }
    ],
    "rules": [
        {
            "type": "field",
            "balancerTag": "cdn-balancer",
            "inboundTag": ["yourfrom"]
        }
    ]
}
```

---

## ECH 更新逻辑说明

本版本程序在运行时会启动一个后台任务：
1.  **轮询机制**：程序会从内置的 Cloudflare 优选 IP 库中轮询，通过 HTTPS 请求 `tlsSettings` 中配置的 `echDohServer`。
2.  **参数透传**：请求时会带上 `domain=ServerName` 参数，确保获取到正确的 ECH 密钥。
3.  **自动应用**：拉取成功后，程序会自动更新全局内存缓存，后续所有经过负载均衡器的 TLS 连接都将使用最新的 ECH 密钥。

---

## 常见问题排查

* **没有生成子出站？**：检查 `tag` 是否严格以 `cdn-` 开头，且 `result/result.json` 路径是否正确。
* **ECH 更新失败？**：请确认 `echDohServer` 填入的是已经在 Cloudflare 绑定了 **自定义域名** 的 Worker 地址，且**不要**带 `https://` 协议头。
* **连接重置 (RST)？**：如果由于 SNI 拦截导致无法更新，程序会自动尝试不同的优选 IP 绕过。

---