package httpfetcher

import (
	"context"
	"io"
	gonet "net"
	"net/http"

	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/net"

	"github.com/gzjjjfree/v5-result/app/subscription"
	"github.com/gzjjjfree/v5-result/app/subscription/documentfetcher"
	"github.com/gzjjjfree/v5-result/common/environment"
	"github.com/gzjjjfree/v5-result/common/environment/envctx"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

func newHTTPFetcher() *httpFetcher {
	return &httpFetcher{}
}

func init() {
	common.Must(documentfetcher.RegisterFetcher("http", newHTTPFetcher()))
}

type httpFetcher struct{}

func (h *httpFetcher) DownloadDocument(ctx context.Context, source *subscription.ImportSource, opts ...documentfetcher.FetcherOptions) ([]byte, error) {
	instanceNetwork := envctx.EnvironmentFromContext(ctx).(environment.InstanceNetworkCapabilitySet)
	outboundDialer := instanceNetwork.OutboundDialer()
	var httpRoundTripper http.RoundTripper //nolint: gosimple
	httpRoundTripper = &http.Transport{
		DialContext: func(ctx_ context.Context, network string, addr string) (gonet.Conn, error) {
			dest, err := net.ParseDestination(network + ":" + addr)
			if err != nil {
				return nil, newError("unable to parse destination")
			}
			return outboundDialer(ctx, dest, source.ImportUsingTag)
		},
	}
	request, err := http.NewRequest("GET", source.Url, nil)
	if err != nil {
		return nil, newError("unable to generate request").Base(err)
	}
	resp, err := httpRoundTripper.RoundTrip(request)
	if err != nil {
		return nil, newError("unable to send request").Base(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, newError("unexpected http status ", resp.StatusCode, "=", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, newError("unable to read response").Base(err)
	}
	return data, nil
}
