package logagent

import (
	"context"
	"log"
	"log-collector/global/setting"
	"log-collector/module/etcd"
	"log-collector/module/logtail"
)

func Run() {
	ctx := context.Background()
	configs, err := etcd.ReadAndWatch(ctx, setting.EtcdSettingCache.LogConfigKey)
	if err != nil {
		log.Println(err)
		return
	}
	// create log reader
	readerMgr, err := logtail.InitAndStart(ctx, configs)
	if err != nil {
		log.Println(err)
		return
	}
	// register observer
	etcd.Register(setting.EtcdSettingCache.LogConfigKey, readerMgr)
}
