package setting

import (
	"log-collector/global/errcode"
	"log-collector/global/globalconst"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	kafkaConfigKey = "Kafka"
	tailConfKey    = "Tail"
	etcdConfKey    = "Etcd"
)

var once sync.Once

var (
	KafkaSettingCache *KafkaSetting
	TailSettingCache  *TailSetting
	EtcdSettingCache  *EtcdSetting
)

func InitSetting(layer int) error {
	var err error
	once.Do(func() {
		vp := viper.New()
		vp.SetConfigName("config")
		vp.AddConfigPath(getAbsPath(layer) + globalconst.PathSeparator + "conf")
		vp.SetConfigType("yaml")
		if e := vp.ReadInConfig(); e != nil {
			err = e
			return
		}
		if e := vp.UnmarshalKey(kafkaConfigKey, &KafkaSettingCache); e != nil {
			err = errcode.InitLogConfigError.WithDetail(e.Error()).ToError()
			return
		}
		if e := vp.UnmarshalKey(tailConfKey, &TailSettingCache); e != nil {
			err = errcode.InitLogConfigError.WithDetail(e.Error()).ToError()
			return
		}
		if e := vp.UnmarshalKey(etcdConfKey, &EtcdSettingCache); e != nil {
			err = errcode.InitLogConfigError.WithDetail(e.Error()).ToError()
			return
		}
	})

	return err
}

func getAbsPath(layer int) string {
	rootPath, _ := os.Getwd()
	for i := 0; i < layer; i++ {
		ind := strings.LastIndex(rootPath, globalconst.PathSeparator)
		rootPath = rootPath[:ind]
	}

	return rootPath
}
