package errcode

const (
	defaultErrorNo = 10000
)

var (
	// 初始化错误区间 [10000, 19999]
	InitKafkaError         = NewApiError(10001, "init kafka error")
	InitLogUtilError       = NewApiError(10002, "init log util error")
	InitLogConfigError     = NewApiError(10003, "init config error")
	InitLogTailReaderError = NewApiError(10004, "init tail reader error")
	// 读取错误区间 [20000, 29999]

)
