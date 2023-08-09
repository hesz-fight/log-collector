package setting

import (
	"log-collector/global/errcode"
	"log-collector/global/globalconst"
	"os"
	"strings"

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

func IntSetting(layer int) error {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath(getAbsPath(layer) + globalconst.PathSeparator + "conf")
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

func getAbsPath(layer int) string {
	rootPath, _ := os.Getwd()
	for i := 0; i < layer; i++ {
		ind := strings.LastIndex(rootPath, globalconst.PathSeparator)
		rootPath = rootPath[:ind]
	}

	return rootPath
}
