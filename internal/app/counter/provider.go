package counter

import (
	"context"
	"fmt"
	"strconv"

	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const increment uint = 100

type Provider interface {
	GetCounter() (*counter, error)
}

type EtcdProvider struct {
	client *etcd.Client
}

func (e *EtcdProvider) GetCounter() (*counter, error) {

	kvc := etcd.NewKV(e.client)

	s, _ := concurrency.NewSession(e.client)
	defer s.Close()
	l := concurrency.NewMutex(s, "/counter-lock/")

	ctx := context.Background()

	if err := l.Lock(ctx); err != nil {
		fmt.Println(err)
	}

	gResp, err := kvc.Get(ctx, "counter")
	if err != nil {
		return nil, err
	}
	counterStart, _ := strconv.Atoi(string(gResp.Kvs[0].Value))
	counterEnd := counterStart + int(increment)

	kvc.Put(ctx, "counter", strconv.Itoa(counterEnd))

	if err := l.Unlock(ctx); err != nil {
		fmt.Println(err)
	}

	return newCounter(uint64(counterStart), uint64(counterEnd)), nil
}

func NewEtcdProvider(client *etcd.Client) *EtcdProvider {
	return &EtcdProvider{client: client}
}
