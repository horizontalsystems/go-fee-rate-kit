package main

import (
    "fmt"
    "github.com/horizontalsystems/go-fee-rate-kit"
)

func main() {
    feeRateKit := feeratekit.NewFeeRateKit(".")

    //feeRateKit.Refresh()

    rate := feeRateKit.Bitcoin()
    fmt.Printf("Bitcoin: low: %v, medium: %v, high: %v\n", rate.Low(), rate.Medium(), rate.High())

    rate = feeRateKit.BitcoinCash()
    fmt.Printf("Bitcoin Cash: low: %v, medium: %v, high: %v\n", rate.Low(), rate.Medium(), rate.High())

    rate = feeRateKit.Dash()
    fmt.Printf("Dash: low: %v, medium: %v, high: %v\n", rate.Low(), rate.Medium(), rate.High())

    rate = feeRateKit.Ethereum()
    fmt.Printf("Ethereum: low: %v, medium: %v, high: %v\n", rate.Low(), rate.Medium(), rate.High())
}
