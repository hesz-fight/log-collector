package main

import (
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/kafka"
	"log-collector/module/logagent"
	"log-collector/module/logtail"
	"time"
)

func main() {
	// init
	if err := setting.IntSetting(); err != nil {
		panic(err)
	}
	if err := kafka.InitProducer(setting.KafkaSettingCache.Addrs); err != nil {
		panic(err)
	}

	if err := etcd.InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}

	// obtain the file path and topic from etcd

	if err := logtail.InitLogReader(setting.TailSettingCache.FilePath,
		setting.TailSettingCache.MaxBufSize); err != nil {
		panic(err)
	}
	// start reading log file
	logagent.StartReadLog()
}
