package logagent

import (
	"log-collector/global/apierror"
	"log-collector/global/setting"
)

var producer *KafkaSyncProducer

func InitProducer() error {
	addrs := setting.KafkaSettingObj.Addrs
	topic := setting.KafkaSettingObj.Topic

	var err error
	producer, err = NewSyncProducer(addrs, topic, nil)
	if err != nil {
		return apierror.InitKafkaError.ToError()
	}

	return nil
}
