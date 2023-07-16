package logreader

import (
	"os"

	"log-collector/global/errcode"

	"github.com/hpcloud/tail"
)

var TailReader *tail.Tail

// 从文件末尾开始读取
func InitLogReader() error {
	filename := "log/tail.log"
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
		MustExist: false,
		Poll:      true,
	}

	var err error
	TailReader, err = tail.TailFile(filename, config)
	if err != nil {
		return errcode.InitLogTailReaderError.ToError()
	}

	return nil
}
