package feeratekit

import (
    "log"
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

func NewFeeRateKit(dataDir string) (*FeeRateKit, error) {
    storage, err := newStorage(dataDir)

    if err != nil {
        return nil, err
    }

    kit := FeeRateKit{
        storage:         storage,
        syncer:          newSyncer(storage),
        refreshChannels: make([]chan interface{}, 0),
    }

    kit.Refresh()

    ticker := time.NewTicker(5 * time.Second)

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
    go func() {
        err := kit.syncer.syncRates()

        if err != nil {
            log.Printf("Refresh failed: %v", err)
            return
        }

        log.Printf("Refresh success")

        for _, channel := range kit.refreshChannels {
            channel <- true
        }
    }()
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
