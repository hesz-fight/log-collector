package main

import (
	"log-collector/global/setting"
	"log-collector/module/logagent"
	"log-collector/module/logreader"
)

func start() {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

}

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

	start()
}
