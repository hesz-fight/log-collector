package logagent

import (
	"log-collector/global/errcode"
	"log-collector/global/setting"
)

var Producer *KafkaSyncProducer

func InitProducer() error {
	addrs := setting.KafkaSettingObj.Addrs
	topic := setting.KafkaSettingObj.Topic

	var err error
	Producer, err = NewSyncProducer(addrs, topic, nil)
	if err != nil {
		return errcode.InitKafkaError.ToError()
	}

	return nil
}
