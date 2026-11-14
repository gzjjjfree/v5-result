package dns_test

import (
	"context"
	"testing"
	"time"

	. "github.com/gzjjjfree/v5-result/app/dns"
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/features/dns"
)

func TestLocalNameServer(t *testing.T) {
	s := NewLocalNameServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ips, err := s.QueryIP(ctx, "google.com", net.IP{}, dns.IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
	}, false)
	cancel()
	common.Must(err)
	if len(ips) == 0 {
		t.Error("expect some ips, but got 0")
	}
}
