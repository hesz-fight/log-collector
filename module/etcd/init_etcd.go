package etcd

import (
	"context"
	"log-collector/global/errcode"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdWrapper struct {
	etcdCli *clientv3.Client
}

func NewEtcdWrapper(endpointds []string, dialTimeout time.Duration) (*EtcdWrapper, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpointds,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, errcode.InitLogEtcdrError.WithDetail(err.Error()).ToError()
	}

	etcdWrapper := &EtcdWrapper{
		etcdCli: cli,
	}

	return etcdWrapper, nil
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

// DeleteOne delete one key
func (e *EtcdWrapper) DeleteOne(ctx context.Context, key string) error {
	_, err := e.etcdCli.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRange delete keys in range [key, end)
func (e *EtcdWrapper) DeleteRange(ctx context.Context, key string, end string) (int64, error) {
	rsp, err := e.etcdCli.Delete(ctx, key, clientv3.WithRange(end))
	if err != nil {
		return 0, nil
	}

	return rsp.Deleted, nil
}
