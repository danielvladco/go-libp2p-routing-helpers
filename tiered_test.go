package routinghelpers

import (
	"context"
	"testing"

	errwrap "github.com/hashicorp/errwrap"
	cid "github.com/ipfs/go-cid"
	routing "github.com/libp2p/go-libp2p-routing"
)

func TestTieredGet(t *testing.T) {
	d := Tiered{
		Routers: []routing.IpfsRouting{
			Null{},
			&Compose{
				ValueStore:     new(dummyValueStore),
				ContentRouting: Null{},
				PeerRouting:    Null{},
			},
			&Compose{
				ValueStore:     new(dummyValueStore),
				ContentRouting: Null{},
				PeerRouting:    Null{},
			},
			&Compose{
				ValueStore:     new(dummyValueStore),
				ContentRouting: Null{},
				PeerRouting:    Null{},
			},
			Null{},
			&Compose{},
		},
	}
	ctx := context.Background()
	if err := d.Routers[1].PutValue(ctx, "k1", []byte("v1")); err != nil {
		t.Fatal(err)
	}
	if err := d.Routers[2].PutValue(ctx, "k2", []byte("v2")); err != nil {
		t.Fatal(err)
	}
	if err := d.Routers[2].PutValue(ctx, "k1", []byte("v1shadow")); err != nil {
		t.Fatal(err)
	}
	if err := d.Routers[3].PutValue(ctx, "k3", []byte("v3")); err != nil {
		t.Fatal(err)
	}

	for k, v := range map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	} {
		actual, err := d.GetValue(ctx, k)
		if err != nil {
			t.Fatal(err)
		}
		if string(actual) != v {
			t.Errorf("expected %s, got %s", v, string(actual))
		}
	}
	if _, err := d.GetValue(ctx, "missing"); err != routing.ErrNotFound {
		t.Fatal("wrong error: ", err)
	}

	if err := d.PutValue(ctx, "key", []byte("value")); err != nil {
		t.Fatal(err)
	}

	if _, err := d.GetValue(ctx, "/error/myErr"); !errwrap.Contains(err, "myErr") {
		t.Fatalf("expected error to contain myErr, got: %s", err)
	}

	if _, err := (Tiered{Routers: []routing.IpfsRouting{d.Routers[1]}}).GetValue(ctx, "/error/myErr"); !errwrap.Contains(err, "myErr") {
		t.Fatalf("expected error to contain myErr, got: %s", err)
	}

	for _, di := range append([]routing.IpfsRouting{d}, d.Routers[1:len(d.Routers)-2]...) {
		v, err := di.GetValue(ctx, "key")
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != "value" {
			t.Errorf("expected value, got %s", string(v))
		}
	}
}

func TestTieredNoSupport(t *testing.T) {
	d := Tiered{Routers: []routing.IpfsRouting{Tiered{Routers: []routing.IpfsRouting{Null{}}}}}
	if _, ok := <-d.FindProvidersAsync(context.Background(), cid.Cid{}, 0); ok {
		t.Fatal("shouldn't have found a provider")
	}
}
