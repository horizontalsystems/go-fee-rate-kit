package feeratekit

import (
    "time"
)

type FeeRateKit struct {
    storage *storage
}

func NewFeeRateKit(dataDir string) *FeeRateKit {
    return &FeeRateKit{
        storage: newStorage(dataDir),
    }
}

func (kit FeeRateKit) feeRate(coin coin) *FeeRate {
    storedFeeRate := kit.storage.feeRate(coin)

    if storedFeeRate != nil {
        return storedFeeRate
    }

    return coin.defaultFeeRate()
}

func (kit FeeRateKit) Bitcoin() *FeeRate {
    return kit.feeRate(bitcoin)
}

func (kit FeeRateKit) BitcoinCash() *FeeRate {
    return kit.feeRate(bitcoinCash)
}

func (kit FeeRateKit) Dash() *FeeRate {
    return kit.feeRate(dash)
}

func (kit FeeRateKit) Ethereum() *FeeRate {
    return kit.feeRate(ethereum)
}

func (kit FeeRateKit) Refresh() {
    kit.storage.save(FeeRate{dash, 100, 200, 300, time.Now().Unix()})
}
