package logtail

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"log-collector/global/errcode"
	"log-collector/model/common"

	"github.com/hpcloud/tail"
)

var TR *TailReader

type TailReader struct {
	// tail chan is bloeked queue
	tail *tail.Tail
	// log buf queue
	buf chan string
	// receive stop signal
	done chan common.Empty
}

func InitLogReader(filename string, bufSize int) error {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
		MustExist: false,
		Poll:      true,
	}
	tail, err := tail.TailFile(filename, config)
	if err != nil {
		return errcode.InitLogTailReaderError.ToError()
	}

	TR = &TailReader{
		tail: tail,
		buf:  make(chan string, bufSize),
		done: make(chan common.Empty),
	}

	return nil
}

func (t *TailReader) Srart() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(string(debug.Stack()))
			}
		}()

		var line *tail.Line
		var ok bool
	LOOP:
		for {
			select {
			case <-t.done:
				break LOOP
			case line, ok = <-t.tail.Lines:
				if !ok {
					// file reopen
					time.Sleep(1 * time.Second)
					continue
				}
				if line.Err != nil {
					// report error
					continue
				}
				t.buf <- line.Text
			}
		}
		fmt.Println("stop read by sigal")
	}()
}

// Done stop read goroutine
func (t *TailReader) Done() chan<- common.Empty {
	return t.done
}

// Reader asynchronous read
func (t *TailReader) AsyncRead() (string, bool) {
	select {
	case r, ok := <-t.buf:
		if !ok {
			return "", false
		}
		return r, true
	default:
		return "", false
	}
}

// SyncRead synchronous read
func (t *TailReader) SyncRead() (string, bool) {
	r, ok := <-t.buf
	if !ok {
		return "", false
	}

	return r, true
}
