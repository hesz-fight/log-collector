package main

import (
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/kafka"
	"log-collector/module/logagent"
	"time"
)

func main() {
	// init
	if err := setting.InitSetting(0); err != nil {
		panic(err)
	}
	if err := kafka.InitProducer(setting.KafkaSettingCache.Addrs); err != nil {
		panic(err)
	}
	if err := etcd.InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}
	// run
	logagent.Run()
}
