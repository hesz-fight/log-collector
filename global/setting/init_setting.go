package setting

import (
	"log-collector/global/errcode"

	"github.com/spf13/viper"
)

const (
	kafkaConfigKey = "kafka"
)

var (
	KafkaSettingCache *KafkaSetting
)

func IntSetting() error {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("conf")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	if err := vp.UnmarshalKey(kafkaConfigKey, KafkaSettingCache); err != nil {
		return errcode.InitLogConfigError.ToError()
	}

	return nil
}
