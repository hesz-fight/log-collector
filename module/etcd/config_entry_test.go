package etcd

import (
	"context"
	"testing"
	"time"
)

func TestConfigEntry(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	dailTimeout := 3 * time.Second
	cli, err := NewEtcdWrapper(endpoints, dailTimeout)
	if err != nil {
		t.Fatalf("NewEtcdWrapper error:%v", err)
	}

	val := `[{"path":"E:/GoModules/log-collector/log/tail.log","topic":"kafka-test-topic"}]`
	if err := cli.Put(context.Background(), "config_test", val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}

	entryArr := []*ConfigEntry{}
	if err := ReadAndUnmarshal(context.Background(), cli, "config_test", &entryArr); err != nil {
		t.Fatalf("ReadAndUnmarshal error:%v", err)
	}

	t.Logf("entryArr:%v", entryArr[0])

	go WatchAndUpdate(context.Background(), cli, "config_test", &entryArr)

	val = `[{"path":"D:/GoModules/log-collector/log/tail.log","topic":"kafka-test-topic"}]`
	if err := cli.Put(context.Background(), "config_test", val); err != nil {
		t.Fatalf("put into etcd error:%v", err)
	}

	time.Sleep(1 * time.Second)

	t.Logf("entryArr: %v", entryArr[0])
}
