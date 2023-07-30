package setting

import (
	"log-collector/global/errcode"

	"github.com/spf13/viper"
)

const (
	kafkaConfigKey = "Kafka"
	tailConfKey    = "Tail"
	etcdConfKey    = "Etcd"
)

var (
	KafkaSettingCache *KafkaSetting
	TailSettingCache  *TailSetting
	EtcdSettingCache  *EtcdSetting
)

func IntSetting() error {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("conf")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	if err := vp.UnmarshalKey(kafkaConfigKey, &KafkaSettingCache); err != nil {
		return errcode.InitLogConfigError.WithDetail(err.Error()).ToError()
	}
	if err := vp.UnmarshalKey(tailConfKey, &TailSettingCache); err != nil {
		return errcode.InitLogConfigError.WithDetail(err.Error()).ToError()
	}
	if err := vp.UnmarshalKey(etcdConfKey, &EtcdSettingCache); err != nil {
		return errcode.InitLogConfigError.WithDetail(err.Error()).ToError()
	}

	return nil
}
