package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log-collector/model/common"
	"reflect"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
)

const (
	maxSize = 1024
)

var (
	defaultEtcdContent *EtcdContent
	once               sync.Once
)

// Observer ...
type Observer interface {
	Notify(context.Context, common.Any)
}

type ConfigEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

type EtcdContent struct {
	Observers map[string][]Observer
	Wrapper   *EtcdWrapper
}

func (e *EtcdContent) register(key string, observer Observer) error {
	if e.Observers == nil {
		e.Observers = make(map[string][]Observer)
	}
	if e.Observers[key] == nil {
		e.Observers[key] = make([]Observer, 0)
	}

	if len(e.Observers[key]) > maxSize {
		return errors.New("get register more observers")
	}

	e.Observers[key] = append(e.Observers[key], observer)
	return nil
}

// notifyAll ...
func (e *EtcdContent) notifyAll(ctx context.Context, key string, conf []*ConfigEntry) {
	log.Println("etcd context notify all obervers for key:" + key)
	observers, ok := e.Observers[key]
	if !ok {
		return
	}
	for _, o := range observers {
		go o.Notify(ctx, conf)
	}
}

// readAndUnmarshal read from etcd and unmarsh
func (e *EtcdContent) readAndUnmarshal(ctx context.Context, key string, val interface{}) error {
	vals, err := e.Wrapper.Get(ctx, key)
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
func (e *EtcdContent) watchAndUpdate(ctx context.Context, key string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("start panic")
		}
	}()

	rspCh := e.Wrapper.Watch(ctx, key)
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
			log.Printf("watching value change key:%v val:%v\n", key, reflect.ValueOf(conf).Index(0).Elem().Interface())
			e.notifyAll(ctx, key, conf)
		}
	}

	return nil
}

// InitEtcd ...
func InitEtcd(endpointds []string, dialTimeout time.Duration) error {
	var err error
	once.Do(
		func() {
			etcdWrapper, e := NewEtcdWrapper(endpointds, dialTimeout)
			if e != nil {
				err = e
				return
			}
			defaultEtcdContent = &EtcdContent{
				Observers: make(map[string][]Observer),
				Wrapper:   etcdWrapper,
			}
		})

	return err
}

// Register ...
func Register(key string, observer Observer) error {
	return defaultEtcdContent.register(key, observer)
}

// ReadAndWatch read and watch the config
func ReadAndWatch(ctx context.Context, configKey string) ([]*ConfigEntry, error) {
	conf := make([]*ConfigEntry, 0)
	if err := defaultEtcdContent.readAndUnmarshal(ctx, configKey, &conf); err != nil {
		return nil, err
	}

	go defaultEtcdContent.watchAndUpdate(ctx, configKey)

	return conf, nil
}
