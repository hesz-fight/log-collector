package main

import (
	"log-collector/global/setting"
	"log-collector/module/logagent"
	"log-collector/module/logreader"
)

func main() {
	if err := setting.IntSetting(); err != nil {
		panic(err)
	}
	if err := logagent.InitProducer(); err != nil {
		panic(err)
	}
	if err := logreader.InitLogReader(); err != nil {
		panic(err)
	}
}
