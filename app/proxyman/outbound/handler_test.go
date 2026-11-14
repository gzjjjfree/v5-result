package outbound_test

import (
	"context"
	"testing"
	_ "unsafe"

	"google.golang.org/protobuf/types/known/anypb"

	core "github.com/gzjjjfree/v5-result"
	"github.com/gzjjjfree/v5-result/app/policy"
	. "github.com/gzjjjfree/v5-result/app/proxyman/outbound"
	"github.com/gzjjjfree/v5-result/app/stats"
	"github.com/gzjjjfree/v5-result/common/environment"
	"github.com/gzjjjfree/v5-result/common/environment/deferredpersistentstorage"
	"github.com/gzjjjfree/v5-result/common/environment/envctx"
	"github.com/gzjjjfree/v5-result/common/environment/filesystemimpl"
	"github.com/gzjjjfree/v5-result/common/environment/systemnetworkimpl"
	"github.com/gzjjjfree/v5-result/common/environment/transientstorageimpl"
	"github.com/gzjjjfree/v5-result/common/net"
	"github.com/gzjjjfree/v5-result/common/serial"
	"github.com/gzjjjfree/v5-result/features/outbound"
	"github.com/gzjjjfree/v5-result/proxy/freedom"
	"github.com/gzjjjfree/v5-result/transport/internet"
)

func TestInterfaces(t *testing.T) {
	_ = (outbound.Handler)(new(Handler))
	_ = (outbound.Manager)(new(Manager))
}

//go:linkname toContext github.com/gzjjjfree/v5-result.toContext
func toContext(ctx context.Context, v *core.Instance) context.Context

func TestOutboundWithoutStatCounter(t *testing.T) {
	config := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&stats.Config{}),
			serial.ToTypedMessage(&policy.Config{
				System: &policy.SystemPolicy{
					Stats: &policy.SystemPolicy_Stats{
						InboundUplink: true,
					},
				},
			}),
		},
	}

	v, _ := core.New(config)
	v.AddFeature((outbound.Manager)(new(Manager)))
	ctx := toContext(context.Background(), v)
	defaultNetworkImpl := systemnetworkimpl.NewSystemNetworkDefault()
	defaultFilesystemImpl := filesystemimpl.NewDefaultFileSystemDefaultImpl()
	deferredPersistentStorageImpl := deferredpersistentstorage.NewDeferredPersistentStorage(ctx)
	rootEnv := environment.NewRootEnvImpl(ctx,
		transientstorageimpl.NewScopedTransientStorageImpl(), defaultNetworkImpl.Dialer(), defaultNetworkImpl.Listener(),
		defaultFilesystemImpl, deferredPersistentStorageImpl)
	proxyEnvironment := rootEnv.ProxyEnvironment("o")
	ctx = envctx.ContextWithEnvironment(ctx, proxyEnvironment)
	h, _ := NewHandler(ctx, &core.OutboundHandlerConfig{
		Tag:           "tag",
		ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
	})
	conn, _ := h.(*Handler).Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146))
	_, ok := conn.(*internet.StatCouterConnection)
	if ok {
		t.Errorf("Expected conn to not be StatCouterConnection")
	}
}

func TestOutboundWithStatCounter(t *testing.T) {
	config := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&stats.Config{}),
			serial.ToTypedMessage(&policy.Config{
				System: &policy.SystemPolicy{
					Stats: &policy.SystemPolicy_Stats{
						OutboundUplink:   true,
						OutboundDownlink: true,
					},
				},
			}),
		},
	}

	v, _ := core.New(config)
	v.AddFeature((outbound.Manager)(new(Manager)))
	ctx := toContext(context.Background(), v)
	defaultNetworkImpl := systemnetworkimpl.NewSystemNetworkDefault()
	defaultFilesystemImpl := filesystemimpl.NewDefaultFileSystemDefaultImpl()
	deferredPersistentStorageImpl := deferredpersistentstorage.NewDeferredPersistentStorage(ctx)
	rootEnv := environment.NewRootEnvImpl(ctx,
		transientstorageimpl.NewScopedTransientStorageImpl(), defaultNetworkImpl.Dialer(), defaultNetworkImpl.Listener(),
		defaultFilesystemImpl, deferredPersistentStorageImpl)
	proxyEnvironment := rootEnv.ProxyEnvironment("o")
	ctx = envctx.ContextWithEnvironment(ctx, proxyEnvironment)
	h, _ := NewHandler(ctx, &core.OutboundHandlerConfig{
		Tag:           "tag",
		ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
	})
	conn, _ := h.(*Handler).Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146))
	_, ok := conn.(*internet.StatCouterConnection)
	if !ok {
		t.Errorf("Expected conn to be StatCouterConnection")
	}
}
