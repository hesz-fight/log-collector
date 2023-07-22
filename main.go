package main

import (
	"log-collector/global/setting"
	"log-collector/module/kafka"
	"log-collector/module/logagent"
	"log-collector/module/logtail"
)

func main() {
	// init
	if err := setting.IntSetting(); err != nil {
		panic(err)
	}
	if err := kafka.InitProducer(); err != nil {
		panic(err)
	}
	if err := logtail.InitLogReader(); err != nil {
		panic(err)
	}

	logagent.StartReadLog()
}
