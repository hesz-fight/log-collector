package main

import (
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/kafka"
	"log-collector/module/logtail"
	"time"
)

func main() {
	// init
	if err := setting.IntSetting(0); err != nil {
		panic(err)
	}
	if err := kafka.InitProducer(setting.KafkaSettingCache.Addrs); err != nil {
		panic(err)
	}
	// init
	if err := etcd.InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}

	configs, err := etcd.ReadAndWatch(setting.EtcdSettingCache.LogConfigKey)
	if err != nil {
		panic(err)
	}
	// create log reader
	rm, err := logtail.InitTailReaders(configs)
	if err != nil {
		panic(err)
	}

	etcd.Register(rm)
	// start reading
	rm.StartReadLog()
}
