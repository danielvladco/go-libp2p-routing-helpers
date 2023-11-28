//go:build notracing

package tracing

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multihash"
)

// Tracer is a noop tracer
type Tracer string

func (t Tracer) StartSpan(ctx context.Context, name string) (context.Context, *noopSpan) {
	return ctx, &noopSpan{}
}

func (t Tracer) Provide(routerName string, ctx context.Context, key cid.Cid, announce bool) (_ context.Context, end func(error)) {
	return ctx, func(err error) {}
}

func (t Tracer) ProvideMany(routerName string, ctx context.Context, keys []multihash.Multihash) (_ context.Context, end func(error)) {
	return ctx, func(err error) {}
}

func (t Tracer) FindProvidersAsync(routerName string, ctx context.Context, key cid.Cid, count int) (_ context.Context, passthrough func(<-chan peer.AddrInfo, error) <-chan peer.AddrInfo) {
	return ctx, func(infos <-chan peer.AddrInfo, err error) <-chan peer.AddrInfo { return infos }
}

func (t Tracer) FindPeer(routerName string, ctx context.Context, id peer.ID) (_ context.Context, end func(peer.AddrInfo, error)) {
	return ctx, func(info peer.AddrInfo, err error) {}
}

func (t Tracer) PutValue(routerName string, ctx context.Context, key string, val []byte, opts ...routing.Option) (_ context.Context, end func(error)) {
	return ctx, func(err error) {}
}

func (t Tracer) GetValue(routerName string, ctx context.Context, key string, opts ...routing.Option) (_ context.Context, end func([]byte, error)) {
	return ctx, func(bytes []byte, err error) {}
}

func (t Tracer) SearchValue(routerName string, ctx context.Context, key string, opts ...routing.Option) (_ context.Context, passthrough func(<-chan []byte, error) (<-chan []byte, error)) {
	return ctx, func(i <-chan []byte, err error) (<-chan []byte, error) { return i, err }
}

func (t Tracer) Bootstrap(routerName string, ctx context.Context) (_ context.Context, end func(error)) {
	return ctx, func(err error) {}
}

type noopSpan struct{}

func (*noopSpan) IsRecording(...any) bool { return false }
func (*noopSpan) SetAttributes(...any)    {}
func (*noopSpan) End(...any)              {}
func (*noopSpan) SetStatus(...any)        {}
func (*noopSpan) AddEvent(...any)         {}
