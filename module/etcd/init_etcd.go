package etcd

import (
	"context"
	"log-collector/global/errcode"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var EtcdWrapperCli *EtcdWrapper

type EtcdWrapper struct {
	etcdCli *clientv3.Client
}

func InitEtcd(endpointds []string, dialTimeout time.Duration) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpointds,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return errcode.InitLogEtcdrError.WithDetail(err.Error()).ToError()
	}

	EtcdWrapperCli = &EtcdWrapper{
		etcdCli: cli,
	}

	return nil
}

// Watch watch the key
func (e *EtcdWrapper) Watch(ctx context.Context, key string) clientv3.WatchChan {
	return e.etcdCli.Watch(ctx, key)
}

// Put ...
func (e *EtcdWrapper) Put(ctx context.Context, key string, val string) error {
	_, err := e.etcdCli.Put(ctx, key, val)
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (e *EtcdWrapper) Get(ctx context.Context, key string) ([]string, error) {
	resp, err := e.etcdCli.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	rsp := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		rsp = append(rsp, string(kv.Value))
	}

	return rsp, nil
}