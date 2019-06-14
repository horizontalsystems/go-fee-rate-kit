package feeratekit

import (
    "time"
)

type FeeRateKit struct {
    storage         *storage
    syncer          *syncer
    refreshChannels []chan interface{}
}

type Handler interface {
    OnRefresh()
}

func NewFeeRateKit(dataDir string, infuraProjectId string, infuraProjectSecret string) (*FeeRateKit, error) {
    storage, err := newStorage(dataDir)

    if err != nil {
        return nil, err
    }

    kit := FeeRateKit{
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

func (kit *FeeRateKit) feeRate(coin string) *FeeRate {
    storedFeeRate := kit.storage.feeRate(coin)

    if storedFeeRate != nil {
        return storedFeeRate
    }

    return defaultFeeRate(coin)
}

func (kit *FeeRateKit) Bitcoin() *FeeRate {
    return kit.feeRate(bitcoin)
}

func (kit *FeeRateKit) BitcoinCash() *FeeRate {
    return kit.feeRate(bitcoinCash)
}

func (kit *FeeRateKit) Dash() *FeeRate {
    return kit.feeRate(dash)
}

func (kit *FeeRateKit) Ethereum() *FeeRate {
    return kit.feeRate(ethereum)
}

func (kit *FeeRateKit) Refresh() {
    kit.syncer.syncRates()
}

func (kit *FeeRateKit) Subscribe(handler Handler) {
    channel := make(chan interface{}, 1)

    kit.refreshChannels = append(kit.refreshChannels, channel)

    go func() {
        for {
            _ = <-channel

            handler.OnRefresh()
        }
    }()
}

func (kit *FeeRateKit) didSyncRates() {
    for _, channel := range kit.refreshChannels {
        channel <- true
    }
}
