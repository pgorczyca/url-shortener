package shortener

import (
	"context"
	"strconv"

	"github.com/pgorczyca/url-shortener/internal/app/utils"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

var c = utils.GetConfig()

type RangeProvider interface {
	GetRange() (*counterRange, error)
	Initialize() error
}

type EtcdRangeProvider struct {
	client *etcd.Client
}

func (e *EtcdRangeProvider) Initialize() error {
	kvc := etcd.NewKV(e.client)

	session, err := concurrency.NewSession(e.client)
	if err != nil {
		utils.Logger.Info("Not able to create initial session", zap.Error(err))
	}
	defer session.Close()

	mutex := concurrency.NewMutex(session, "/counter-lock/")
	ctx := context.Background()

	if err := mutex.Lock(ctx); err != nil {
		utils.Logger.Error("Not able to lock mutex.", zap.Error(err))
	}

	gResp, err := kvc.Get(ctx, "lastCounterEnd")
	if err != nil {
		utils.Logger.Error("Not able to get lastCounterEnd from etcd.", zap.Error(err))
		return err
	}

	if len(gResp.Kvs) == 0 {
		kvc.Put(ctx, "lastCounterEnd", "0")
	}

	if err := mutex.Unlock(ctx); err != nil {
		utils.Logger.Error("Not able to unlock mutex.", zap.Error(err))
	}
	return nil
}

func (e *EtcdRangeProvider) GetRange() (*counterRange, error) {

	kvc := etcd.NewKV(e.client)

	session, err := concurrency.NewSession(e.client)
	defer session.Close()
	mutex := concurrency.NewMutex(session, "/counter-lock/")

	ctx := context.Background()

	if err := mutex.Lock(ctx); err != nil {
		utils.Logger.Error("Not able to lock mutex.", zap.Error(err))
	}

	gResp, err := kvc.Get(ctx, "lastCounterEnd")
	if err != nil {
		utils.Logger.Error("Not able to get lastCounterEnd from etcd.", zap.Error(err))
		return nil, err
	}
	counterStart, err := strconv.Atoi(string(gResp.Kvs[0].Value))
	counterEnd := counterStart + int(c.CounterIncrement)

	kvc.Put(ctx, "counter", strconv.Itoa(counterEnd))

	if err := mutex.Unlock(ctx); err != nil {
		utils.Logger.Error("Not able to unlock mutex.", zap.Error(err))
	}
	counterTreshold := uint64(counterStart) + uint64((float64(c.CounterIncrement) * c.CounterTreshold))
	return &counterRange{start: uint64(counterStart), end: uint64(counterEnd), treshold: counterTreshold}, nil
}

func NewEtcdProvider(client *etcd.Client) *EtcdRangeProvider {
	return &EtcdRangeProvider{client: client}
}
