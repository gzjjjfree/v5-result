package hysteria2_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/buf"
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/common/protocol/tls/cert"
	"github.com/gzjjjfree/v5-result/common/session"
	"github.com/gzjjjfree/v5-result/testing/servers/udp"
	"github.com/gzjjjfree/v5-result/transport/internet"
	"github.com/gzjjjfree/v5-result/transport/internet/hysteria2"
	"github.com/gzjjjfree/v5-result/transport/internet/tls"
)

func TestTCP(t *testing.T) {
	port := udp.PickPort()

	listener, err := hysteria2.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "hysteria2",
		ProtocolSettings: &hysteria2.Config{Password: "123"},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			Certificate: []*tls.Certificate{
				tls.ParseCertificate(
					cert.MustGenerate(nil,
						cert.DNSNames("www.v2fly.org"),
					),
				),
			},
		},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					fmt.Println(err)
					return
				}
				common.Must2(conn.Write(b.Bytes()))
			}
		}()
	})
	common.Must(err)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := hysteria2.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "hysteria2",
		ProtocolSettings: &hysteria2.Config{Password: "123"},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			ServerName:    "www.v2fly.org",
			AllowInsecure: true,
		},
	})
	common.Must(err)
	defer conn.Close()

	const N = 1000
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}
}

func TestUDP(t *testing.T) {
	port := udp.PickPort()

	listener, err := hysteria2.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "hysteria2",
		ProtocolSettings: &hysteria2.Config{Password: "123", UseUdpExtension: true},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			Certificate: []*tls.Certificate{
				tls.ParseCertificate(
					cert.MustGenerate(nil,
						cert.DNSNames("www.v2fly.org"),
					),
				),
			},
		},
	}, func(conn internet.Connection) {
		fmt.Println("incoming")
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					fmt.Println(err)
					return
				}
				common.Must2(conn.Write(b.Bytes()))
			}
		}()
	})
	common.Must(err)

	defer listener.Close()

	time.Sleep(time.Second)

	address, err := net.ParseDestination("udp:127.0.0.1:1180")
	common.Must(err)
	dctx := session.ContextWithOutbound(context.Background(), &session.Outbound{Target: address})

	conn, err := hysteria2.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "hysteria2",
		ProtocolSettings: &hysteria2.Config{Password: "123", UseUdpExtension: true},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			ServerName:    "www.v2fly.org",
			AllowInsecure: true,
		},
	})
	common.Must(err)
	defer conn.Close()

	const N = 1000
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	common.Must2(conn.Write(b1))

	b2 := buf.New()
	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}
}
