package etcd

import (
	"context"
	"fmt"
	"log-collector/global/ip"
	"log-collector/global/setting"
	"testing"
	"time"
)

func TestConfigEntry(t *testing.T) {
	if err := setting.InitSetting(2); err != nil {
		panic(err)
	}
	// init
	if err := InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}
	// get local ip
	localIP, err := ip.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	key := fmt.Sprintf(setting.EtcdSettingCache.LogConfigKey, localIP)
	val := `[{"path":"E:/GoModules/log-collector/log/tail.log","topic":"kafka-test-topic"}]`
	if err := defaultEtcdContent.Wrapper.Put(context.Background(), key, val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}

	defer func() {
		defaultEtcdContent.Wrapper.DeleteOne(context.Background(), key)
	}()

	ctx := context.Background()

	entryArr, err := ReadAndWatch(ctx, key)
	if err != nil {
		t.Fatalf("ReadAndUnmarshal error:%v", err)
	}
	t.Logf("entryArr:%v", entryArr[0])

	time.Sleep(1 * time.Second)

	val = `[{"path":"E:/GoModules/log-collector/log/tail_test.log","topic":"kafka-test-topic"}]`
	if err := defaultEtcdContent.Wrapper.Put(context.Background(), key, val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}

	time.Sleep(30 * 60 * time.Second)
}
