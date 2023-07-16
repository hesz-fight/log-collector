package setting

import (
	"log-collector/global/apierror"

	"github.com/spf13/viper"
)

const (
	kafkaConfigKey = "kafka"
)

var (
	KafkaSettingObj *KafkaSetting
)

func IntSetting() error {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("conf")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	if err := vp.UnmarshalKey(kafkaConfigKey, KafkaSettingObj); err != nil {
		return apierror.InitLogConfigError.ToError()
	}

	return nil
}
