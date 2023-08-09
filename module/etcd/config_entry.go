package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log-collector/global/setting"
	"reflect"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
)

const (
	maxSize = 1024
)

var defaultEtcdContent *EtcdContent

type Observer interface {
	Notify(conf []*ConfigEntry, bufSize int)
}

type ConfigEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

type EtcdContent struct {
	Ctx       context.Context
	Observers []Observer
	Wrapper   *EtcdWrapper
}

func (e *EtcdContent) Register(observer Observer) error {
	if e.Observers == nil {
		e.Observers = make([]Observer, 0)
	}
	if len(e.Observers) > maxSize {
		return errors.New("get register more observers")
	}

	e.Observers = append(e.Observers, observer)
	return nil
}

func (e *EtcdContent) NotifyAll(conf []*ConfigEntry, bufSize int) {
	for _, o := range e.Observers {
		o.Notify(conf, bufSize)
	}
}

// readAndUnmarshal read from etcd and unmarsh
func (e *EtcdContent) readAndUnmarshal(key string, val interface{}) error {
	vals, err := e.Wrapper.Get(e.Ctx, key)
	if err != nil {
		return err
	}
	for _, v := range vals {
		if err := json.Unmarshal([]byte(v), val); err != nil {
			return err
		}
	}

	return nil
}

// watchAndUpdate watch and update the val
func (e *EtcdContent) watchAndUpdate(key string) error {
	rspCh := e.Wrapper.Watch(e.Ctx, key)
	for rsp := range rspCh {
		for _, event := range rsp.Events {
			var v []byte
			if event.Type == mvccpb.PUT {
				v = event.Kv.Value
			} else if event.Type == mvccpb.DELETE {
				v = []byte{}
			}
			conf := make([]*ConfigEntry, 0)
			if err := json.Unmarshal(v, &conf); err != nil {
				return err
			}
			// notify
			e.NotifyAll(conf, setting.TailSettingCache.MaxBufSize)
			log.Printf("watching value change key:%v val:%v\n", key, reflect.ValueOf(conf).Elem().Index(0).Elem().Interface())
		}
	}

	return nil
}

// Register ...
func Register(observer Observer) error {
	return defaultEtcdContent.Register(observer)
}

// InitEtcd ...
func InitEtcd(endpointds []string, dialTimeout time.Duration) error {
	etcdWrapper, err := NewEtcdWrapper(endpointds, dialTimeout)
	if err != nil {
		return err
	}

	defaultEtcdContent = &EtcdContent{
		Ctx:       context.Background(),
		Observers: make([]Observer, 0),
		Wrapper:   etcdWrapper,
	}

	return nil
}

// ReadAndWatch ...
func ReadAndWatch(configKey string) ([]*ConfigEntry, error) {
	conf := make([]*ConfigEntry, 0)
	if err := defaultEtcdContent.readAndUnmarshal(configKey, &conf); err != nil {
		return nil, err
	}

	go defaultEtcdContent.watchAndUpdate(configKey)

	return conf, nil
}
