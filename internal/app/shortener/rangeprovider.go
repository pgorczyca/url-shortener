package shortener

import (
	"context"
	"fmt"
	"strconv"

	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const increment uint64 = 100

type RangeProvider interface {
	GetRange() (*counterRange, error)
}

type EtcdRangeProvider struct {
	client *etcd.Client
}

func (e *EtcdRangeProvider) GetRange() (*counterRange, error) {

	kvc := etcd.NewKV(e.client)

	session, _ := concurrency.NewSession(e.client)
	defer session.Close()
	mutex := concurrency.NewMutex(session, "/counter-lock/")

	ctx := context.Background()

	if err := mutex.Lock(ctx); err != nil {
		fmt.Println(err)
	}

	gResp, err := kvc.Get(ctx, "counter")
	if err != nil {
		return nil, err
	}
	counterStart, _ := strconv.Atoi(string(gResp.Kvs[0].Value))
	counterEnd := counterStart + int(increment)

	kvc.Put(ctx, "counter", strconv.Itoa(counterEnd))

	if err := mutex.Unlock(ctx); err != nil {
		fmt.Println(err)
	}
	counterTreshold := uint64(counterStart) + uint64((float64(increment) * 0.9))
	return &counterRange{start: uint64(counterStart), end: uint64(counterEnd), treshold: counterTreshold}, nil
}

func NewEtcdProvider(client *etcd.Client) *EtcdRangeProvider {
	return &EtcdRangeProvider{client: client}
}
