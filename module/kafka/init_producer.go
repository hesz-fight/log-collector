package kafka

import (
	"log-collector/global/errcode"
	"log-collector/global/setting"
)

var Producer *KafkaSyncProducer

func InitProducer() error {
	addrs := setting.KafkaSettingCache.Addrs
	topic := setting.KafkaSettingCache.Topic

	var err error
	Producer, err = NewSyncProducer(addrs, topic, nil)
	if err != nil {
		return errcode.InitKafkaError.WithDetail(err.Error()).ToError()
	}

	return nil
}
