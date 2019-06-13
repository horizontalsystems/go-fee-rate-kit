package main

import (
    "fmt"
    "github.com/horizontalsystems/go-fee-rate-kit"
    "log"
)

func main() {
    feeRateKit, err := feeratekit.NewFeeRateKit(".")

    if err != nil {
        log.Fatalln(err)
    }

    logRate(feeRateKit.Bitcoin(), "Bitcoin")
    logRate(feeRateKit.BitcoinCash(), "Bitcoin Cash")
    logRate(feeRateKit.Dash(), "Dash")
    logRate(feeRateKit.Ethereum(), "Ethereum")

    feeRateKit.Subscribe(HandlerStub(func() {
        log.Println("On refresh First")
    }))

    feeRateKit.Subscribe(HandlerStub(func() {
        log.Println("On refresh Second")
    }))

    _, _ = fmt.Scanln()
}

func logRate(rate *feeratekit.FeeRate, name string) {
    log.Printf("%v: low: %v, medium: %v, high: %v\n", name, rate.Low(), rate.Medium(), rate.High())
}

type HandlerStub func()

func (stub HandlerStub) OnRefresh() {
    stub()
}
