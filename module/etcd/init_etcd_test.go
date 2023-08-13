package etcd

import (
	"context"
	"fmt"
	"log-collector/global/ip"
	"log-collector/global/setting"
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

func TestGetFromEtcd(t *testing.T) {
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
	if len(rst) > 0 {
		t.Logf("v:%v", rst[0])
	}
}

func TestPutFromEtcd(t *testing.T) {
	if err := setting.InitSetting(2); err != nil {
		panic(err)
	}

	cli, err := NewEtcdWrapper(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout))
	if err != nil {
		t.Fatalf("NewEtcdWrapper error:%v", err)
	}
	// get local ip
	localIP, err := ip.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	key := fmt.Sprintf(setting.EtcdSettingCache.LogConfigKey, localIP)
	val := `[{"path":"E:/GoModules/log-collector/log/test.log","topic":"kafka-test-topic"}]`
	ctx := context.Background()
	if err := cli.Put(ctx, key, val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}
	// test get
	rst, err := cli.Get(ctx, key)
	if err != nil {
		t.Fatalf("read from etcd error:%v", err)
	}
	t.Logf("len:%v", len(rst))
	if len(rst) > 0 {
		t.Logf("k:%v v:%v", key, rst[0])
	}
}
