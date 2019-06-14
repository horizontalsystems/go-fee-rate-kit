package feeratekit

import (
    "time"
)

type Kit struct {
    storage         *storage
    syncer          *syncer
    refreshChannels []chan interface{}
}

type Handler interface {
    OnRefresh()
}

func NewKit(dataDir string, infuraProjectId string, infuraProjectSecret string) (*Kit, error) {
    storage, err := newStorage(dataDir)

    if err != nil {
        return nil, err
    }

    kit := Kit{
        storage:         storage,
        syncer:          newSyncer(storage, infuraProjectId, infuraProjectSecret),
        refreshChannels: make([]chan interface{}, 0),
    }

    kit.syncer.delegate = &kit

    kit.Refresh()

    ticker := time.NewTicker(30 * time.Second)

    go func() {
        for range ticker.C {
            kit.Refresh()
        }
    }()

    return &kit, nil
}

func (kit *Kit) feeRate(coin string) *FeeRate {
    storedFeeRate := kit.storage.feeRate(coin)

    if storedFeeRate != nil {
        return storedFeeRate
    }

    return defaultFeeRate(coin)
}

func (kit *Kit) Bitcoin() *FeeRate {
    return kit.feeRate(bitcoin)
}

func (kit *Kit) BitcoinCash() *FeeRate {
    return kit.feeRate(bitcoinCash)
}

func (kit *Kit) Dash() *FeeRate {
    return kit.feeRate(dash)
}

func (kit *Kit) Ethereum() *FeeRate {
    return kit.feeRate(ethereum)
}

func (kit *Kit) Refresh() {
    kit.syncer.syncRates()
}

func (kit *Kit) Subscribe(handler Handler) {
    channel := make(chan interface{}, 1)

    kit.refreshChannels = append(kit.refreshChannels, channel)

    go func() {
        for {
            _ = <-channel

            handler.OnRefresh()
        }
    }()
}

func (kit *Kit) didSyncRates() {
    for _, channel := range kit.refreshChannels {
        channel <- true
    }
}
