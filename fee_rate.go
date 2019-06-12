package feeratekit

type FeeRate struct {
    coin           coin
    lowPriority    int64
    mediumPriority int64
    highPriority   int64
    timestamp      int64
}

func (feeRate FeeRate) Low() int64 {
    return max(min(feeRate.lowPriority, feeRate.coin.maxFee()), feeRate.coin.minFee())
}

func (feeRate FeeRate) Medium() int64 {
    return max(min(feeRate.mediumPriority, feeRate.coin.maxFee()), feeRate.coin.minFee())
}

func (feeRate FeeRate) High() int64 {
    return max(min(feeRate.highPriority, feeRate.coin.maxFee()), feeRate.coin.minFee())
}
