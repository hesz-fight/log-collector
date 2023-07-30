package errcode

const (
	defaultErrorNo = 9999
)

var (
	// init error in [10000, 19999]
	InitKafkaError         = NewWrapError(10001, "init kafka error")
	InitLogUtilError       = NewWrapError(10002, "init log util error")
	InitLogConfigError     = NewWrapError(10003, "init config error")
	InitLogTailReaderError = NewWrapError(10004, "init tail reader error")
	InitLogEtcdrError      = NewWrapError(10005, "init etcd reader error")
	// read error in [20000, 29999]

)
