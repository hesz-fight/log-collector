package setting

type KafkaSetting struct {
	Addrs []string
}

type TailSetting struct {
	FilePath   string
	MaxBufSize int
}

type EtcdSetting struct {
	Endpoints    []string
	DialTimeout  int64
	LogConfigKey string
}
