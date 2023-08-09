package logtail

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"log-collector/global/errcode"
	"log-collector/global/setting"
	"log-collector/model/common"
	"log-collector/module/etcd"
	"log-collector/module/kafka"

	"github.com/hpcloud/tail"
)

const (
	retryTime = 3
)

type TailReaderManager struct {
	TailReaders []*TailReader
	Mut         sync.RWMutex
}

type TailReader struct {
	// tail chan is bloeked queue
	tail *tail.Tail
	// log buf queue
	buf chan string
	// receive stop signal
	done chan common.Empty
	// log path
	path string
	// kafka topics
	topic []string
}

// InitTailReaders ...
func InitTailReaders(configEntries []*etcd.ConfigEntry) (*TailReaderManager, error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
		MustExist: false,
		Poll:      true,
	}

	readers := make([]*TailReader, 0, len(configEntries))
	for _, entry := range configEntries {
		tail, err := tail.TailFile(entry.Path, config)
		if err != nil {
			return nil, errcode.InitLogTailReaderError.ToError()
		}
		readers = append(readers, &TailReader{
			tail:  tail,
			buf:   make(chan string, setting.TailSettingCache.MaxBufSize),
			done:  make(chan common.Empty),
			path:  entry.Path,
			topic: strings.Split(entry.Topic, ","),
		})
	}

	readerManager := &TailReaderManager{
		TailReaders: readers,
		Mut:         sync.RWMutex{},
	}

	return readerManager, nil
}

func (t *TailReaderManager) Notify(configEntries []*etcd.ConfigEntry, bufSize int) {
	t.Mut.Lock()
	defer t.Mut.RUnlock()
	// update and done
	newReaders := make([]*TailReader, 0, len(configEntries))
	for _, entry := range configEntries {
		tail, err := tail.TailFile(entry.Path,
			tail.Config{
				ReOpen:    true,
				Follow:    true,
				Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
				MustExist: false,
				Poll:      true,
			})
		if err != nil {
			log.Println(errcode.InitLogTailReaderError.ToError())
			return
		}
		newReaders = append(newReaders, &TailReader{
			tail:  tail,
			buf:   make(chan string, bufSize),
			done:  make(chan common.Empty),
			path:  entry.Path,
			topic: strings.Split(entry.Topic, ","),
		})
	}

	oldReaders := t.TailReaders
	t.TailReaders = newReaders
	t.StartReadLog()
	for _, r := range oldReaders {
		close(r.Done())
	}
}

func (t *TailReaderManager) StartReadLog() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("start panic")
		}
	}()

	fmt.Println("start log agent")
	// parallel read
	for _, reader := range t.TailReaders {
		go func(r *TailReader) {
			r.Srart()
			for {
				text, ok := r.SyncRead()
				if !ok {
					continue
				}
				// send log text to kafka
				for i := 0; i < retryTime; i++ {
					partition, offset, err := kafka.Producer.SendMessag("", text)
					if err != nil {
						continue
					}
					fmt.Printf("send kafka successfully. partition:%v offset:%v\n", partition, offset)
					break
				}
			}

		}(reader)
	}
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
