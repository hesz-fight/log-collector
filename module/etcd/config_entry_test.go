package etcd

import (
	"context"
	"log-collector/global/setting"
	"testing"
	"time"
)

func TestConfigEntry(t *testing.T) {
	if err := setting.IntSetting(2); err != nil {
		panic(err)
	}
	// init
	if err := InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}

	key := setting.EtcdSettingCache.LogConfigKey
	val := `[{"path":"E:/GoModules/log-collector/log/tail.log","topic":"kafka-test-topic"}]`
	if err := defaultEtcdContent.Wrapper.Put(context.Background(), key, val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}

	defer func() {
		defaultEtcdContent.Wrapper.DeleteOne(context.Background(), key)
	}()

	entryArr, err := ReadAndWatch(key)
	if err != nil {
		t.Fatalf("ReadAndUnmarshal error:%v", err)
	}

	t.Logf("entryArr:%v", entryArr[0])
}
