package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
)

const (
	configKey = "config_entry"
)

var configEntries []*ConfigEntry

type ConfigEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func InitEtcd(endpointds []string, dialTimeout time.Duration) error {
	c, err := NewEtcdWrapper(endpointds, dialTimeout)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := ReadAndUnmarshal(ctx, c, configKey, &configEntries); err != nil {
		return err
	}

	go WatchAndUpdate(ctx, c, configKey, &configEntries)

	return nil
}

// ReadAndUnmarshal read from etcd and unmarsh
func ReadAndUnmarshal(ctx context.Context, c *EtcdWrapper, key string, val interface{}) error {
	vals, err := c.Get(ctx, key)
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

// WatchAndUpdate watch and update the val
func WatchAndUpdate(ctx context.Context, c *EtcdWrapper, key string, val interface{}) error {
	// val must be pointer
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return errors.New("val must be pointer")
	}
	rspCh := c.Watch(ctx, key)
	for rsp := range rspCh {
		for _, event := range rsp.Events {
			var v []byte
			if event.Type == mvccpb.PUT {
				v = event.Kv.Value
			} else if event.Type == mvccpb.DELETE {
				v = []byte{}
			}
			if err := json.Unmarshal(v, val); err != nil {
				return err
			}
			log.Printf("watching value change key:%v val:%v\n", key, reflect.ValueOf(val).Elem().Index(0).Elem().Interface())
		}
	}

	return nil
}
