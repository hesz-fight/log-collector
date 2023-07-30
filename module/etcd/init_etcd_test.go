package etcd

import (
	"context"
	"testing"
	"time"
)

func TestEtcdWrapper(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	dailTimeout := 3 * time.Second
	cli, err := NewEtcdWrapper(endpoints, dailTimeout)
	if err != nil {
		t.Fatalf("NewEtcdWrapper error:%v", err)
	}
	ctx := context.Background()
	// test get
	rst, err := cli.Get(ctx, "test")
	if err != nil {
		t.Fatalf("read from etcd error:%v", err)
	}
	t.Logf("len:%v", len(rst))
	t.Logf("v:%v", rst[0])
	// test put
	if err := cli.Put(ctx, "test", "world"); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}
	// test get
	rst, err = cli.Get(ctx, "test")
	if err != nil {
		t.Fatalf("read from etcd error:%v", err)
	}
	t.Logf("len:%v", len(rst))
	t.Logf("v:%v", rst[0])
	if err := cli.Put(ctx, "test", "hello"); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}
}
