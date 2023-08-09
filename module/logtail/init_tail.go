package logtail

import (
	"context"
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
	TailWorker []*TailReader
	Mtx        sync.RWMutex
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
	topics []string
}

// InitAndStart ...
func InitAndStart(ctx context.Context, configEntries []*etcd.ConfigEntry) (*TailReaderManager, error) {
	readers := make([]*TailReader, 0, len(configEntries))
	for _, entry := range configEntries {
		tail, err := tail.TailFile(entry.Path, tail.Config{
			ReOpen:    true,
			Follow:    true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
			MustExist: false,
			Poll:      true,
		})
		if err != nil {
			return nil, errcode.InitLogTailReaderError.ToError()
		}
		readers = append(readers, &TailReader{
			tail:   tail,
			buf:    make(chan string, setting.TailSettingCache.MaxBufSize),
			done:   make(chan common.Empty),
			path:   entry.Path,
			topics: strings.Split(entry.Topic, ","),
		})
	}

	mgr := &TailReaderManager{
		TailWorker: readers,
		Mtx:        sync.RWMutex{},
	}

	mgr.startReadLog(ctx)

	return mgr, nil
}

func (t *TailReaderManager) Notify(ctx context.Context, conf common.Any) {
	t.Mtx.Lock()
	defer t.Mtx.RUnlock()
	confEntryArr, ok := conf.([]*etcd.ConfigEntry)
	if !ok {
		return
	}
	// update and done
	newReaders := make([]*TailReader, 0, len(confEntryArr))
	for _, entry := range confEntryArr {
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
			tail:   tail,
			buf:    make(chan string, setting.TailSettingCache.MaxBufSize),
			done:   make(chan common.Empty),
			path:   entry.Path,
			topics: strings.Split(entry.Topic, ","),
		})
	}

	oldReaders := t.TailWorker
	t.TailWorker = newReaders

	t.startReadLog(ctx)
	for _, r := range oldReaders {
		close(r.Done())
	}
}

func (t *TailReaderManager) startReadLog(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("start panic")
		}
	}()
	log.Println("start log agent")
	// parallel read
	for _, reader := range t.TailWorker {
		go func(r *TailReader) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(string(debug.Stack()))
				}
			}()
			r.Start()
			for {
				text, ok := r.SyncRead()
				if !ok {
					break
				}
				// send log text to kafka
				for _, topic := range r.topics {
					for i := 0; i < retryTime; i++ {
						partition, offset, err := kafka.Producer.SendMessag(topic, text)
						if err != nil {
							continue
						}
						fmt.Printf("send kafka successfully. partition:%v offset:%v\n", partition, offset)
						break
					}
				}
			}
		}(reader)
	}
}

func (t *TailReader) Start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(string(debug.Stack()))
			}
			// close the chan when exiting the func
			close(t.buf)
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
		log.Println("stop read by sigal")
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
