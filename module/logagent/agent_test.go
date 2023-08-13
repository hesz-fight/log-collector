package logagent

import (
	"context"
	"fmt"
	"log"
	"log-collector/global/ip"
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/kafka"
	"log-collector/module/logtail"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// init
	if err := setting.InitSetting(2); err != nil {
		panic(err)
	}
	if err := kafka.InitProducer(setting.KafkaSettingCache.Addrs); err != nil {
		panic(err)
	}
	if err := etcd.InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}

	ctx := context.Background()
	t.Run("", func(t *testing.T) {
		// get local ip
		localIP, err := ip.GetOutBoundIP()
		if err != nil {
			panic(err)
		}
		key := fmt.Sprintf(setting.EtcdSettingCache.LogConfigKey, localIP)
		configs, e := etcd.ReadAndWatch(ctx, key)
		if e != nil {
			t.Fatalf("read and watch fail. err:%v", e)
		}
		if len(configs) == 0 {
			return
		}
		mgr, err := logtail.InitAndStart(ctx, configs)
		if err != nil {
			log.Println(err)
			t.Fatalf("init and start fail. err:%v", e)
			return
		}
		// register observer
		etcd.Register(key, mgr)

		for {
			log.Println("log agent is runing...")
			time.Sleep(1 * 24 * 60 * 60 * time.Second)
		}
	})
}
