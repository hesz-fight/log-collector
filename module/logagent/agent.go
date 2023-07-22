package logagent

import (
	"fmt"
	"log-collector/module/kafka"
	"log-collector/module/logtail"
)

const (
	retryTime = 3
)

func StartReadLog() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("start panic")
		}
	}()

	fmt.Println("start log agent")

	logtail.TR.Srart()
	for {
		text, ok := logtail.TR.SyncRead()
		if !ok {
			continue
		}
		// send log text to kafka
		for i := 0; i < retryTime; i++ {
			partition, offset, err := kafka.Producer.SendMessag(text)
			if err != nil {
				continue
			}
			fmt.Printf("send kafka successfully. partition:%v offset:%v\n", partition, offset)
			break
		}
	}
}
