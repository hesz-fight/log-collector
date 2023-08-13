package main

import (
	"context"
	"fmt"
	"log-collector/global/ip"
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/kafka"
	"log-collector/module/logagent"
	"time"
)

func main() {
	// init setting
	if err := setting.InitSetting(0); err != nil {
		panic(err)
	}
	// init kafka
	if err := kafka.InitProducer(setting.KafkaSettingCache.Addrs); err != nil {
		panic(err)
	}
	// init etcd
	if err := etcd.InitEtcd(setting.EtcdSettingCache.Endpoints,
		time.Duration(setting.EtcdSettingCache.DialTimeout)*time.Second); err != nil {
		panic(err)
	}
	// get local ip
	localIP, err := ip.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	logAgentName2Key := map[string]string{
		"Agent_1": fmt.Sprintf(setting.EtcdSettingCache.LogConfigKey, localIP),
	}
	// run one agent for one key
	for name, key := range logAgentName2Key {
		logagent.NewLogAgent(name).Run(context.Background(), key)
	}
}
