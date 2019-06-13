package feeratekit

const (
    bitcoin     = "BTC"
    bitcoinCash = "BCH"
    dash        = "DASH"
    ethereum    = "ETH"
)

func defaultFeeRate(coin string) *FeeRate {
    switch coin {
    case bitcoin:
        return &FeeRate{coin, 20, 40, 80, 1543211299}
    case bitcoinCash:
        return &FeeRate{coin, 1, 3, 5, 1543211299}
    case dash:
        return &FeeRate{coin, 1, 1, 2, 1557224025}
    case ethereum:
        return &FeeRate{coin, 13000000000, 16000000000, 19000000000, 1543211299}
    }

    return nil
}

func maxFee(coin string) int64 {
    switch coin {
    case bitcoin:
        return 5000
    case bitcoinCash:
        return 500
    case dash:
        return 500
    case ethereum:
        return 3000000000000
    }

    return 0
}

func minFee(coin string) int64 {
    switch coin {
    case bitcoin:
        return 1
    case bitcoinCash:
        return 1
    case dash:
        return 1
    case ethereum:
        return 100000000
    }

    return 0
}
