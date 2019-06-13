package feeratekit

type FeeRate struct {
    coin           string
    lowPriority    int64
    mediumPriority int64
    highPriority   int64
    timestamp      int64
}

func (feeRate *FeeRate) limitedValue(value int64) int64 {
    return max(min(value, maxFee(feeRate.coin)), minFee(feeRate.coin))
}

func (feeRate *FeeRate) Low() int64 {
    return feeRate.limitedValue(feeRate.lowPriority)
}

func (feeRate *FeeRate) Medium() int64 {
    return feeRate.limitedValue(feeRate.mediumPriority)
}

func (feeRate *FeeRate) High() int64 {
    return feeRate.limitedValue(feeRate.highPriority)
}
