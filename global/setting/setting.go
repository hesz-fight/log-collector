package setting

type KafkaSetting struct {
	Addrs []string
	Topic string
}

type TailSetting struct {
	FilePath   string
	MaxBufSize int
}
