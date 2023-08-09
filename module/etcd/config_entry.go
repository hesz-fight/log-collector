package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log-collector/model/common"
	"reflect"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
)

const (
	maxSize = 1024
)

var defaultEtcdContent *EtcdContent

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

func (e *EtcdContent) Register(key string, observer Observer) error {
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

// NotifyAll ...
func (e *EtcdContent) NotifyAll(ctx context.Context, key string, conf []*ConfigEntry) {
	observers, ok := e.Observers[key]
	if !ok {
		return
	}
	for _, o := range observers {
		o.Notify(ctx, conf)
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
			e.NotifyAll(ctx, key, conf)
			log.Printf("watching value change key:%v val:%v\n", key, reflect.ValueOf(conf).Index(0).Elem().Interface())
		}
	}

	return nil
}

// Register ...
func Register(key string, observer Observer) error {
	return defaultEtcdContent.Register(key, observer)
}

// InitEtcd ...
func InitEtcd(endpointds []string, dialTimeout time.Duration) error {
	etcdWrapper, err := NewEtcdWrapper(endpointds, dialTimeout)
	if err != nil {
		return err
	}

	defaultEtcdContent = &EtcdContent{
		Observers: make(map[string][]Observer),
		Wrapper:   etcdWrapper,
	}

	return nil
}

// ReadAndWatch ...
func ReadAndWatch(ctx context.Context, configKey string) ([]*ConfigEntry, error) {
	conf := make([]*ConfigEntry, 0)
	if err := defaultEtcdContent.readAndUnmarshal(ctx, configKey, &conf); err != nil {
		return nil, err
	}

	go defaultEtcdContent.watchAndUpdate(ctx, configKey)

	return conf, nil
}
