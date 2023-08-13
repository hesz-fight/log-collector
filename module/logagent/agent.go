package logagent

import (
	"context"
	"log"
	"log-collector/model/common"
	"log-collector/module/etcd"
	"log-collector/module/logtail"
	"time"
)

type LogAgent struct {
	name string
	exit chan common.Empty
}

func (a *LogAgent) Close() {
	close(a.exit)
}

func (a *LogAgent) Run(ctx context.Context, configKey string) {
	configs, err := etcd.ReadAndWatch(ctx, configKey)
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
	etcd.Register(configKey, readerMgr)

	a.block()
}

func (a *LogAgent) block() {
	select {
	case <-a.exit:
	default:
		log.Println("log agent is running...")
		time.Sleep(1 * 24 * 60 * 60 * time.Second)
	}
}

func NewLogAgent(name string) *LogAgent {
	return &LogAgent{
		name: name,
		exit: make(chan common.Empty),
	}
}
