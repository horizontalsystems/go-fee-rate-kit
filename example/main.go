package main

import (
    "fmt"
    "github.com/horizontalsystems/go-fee-rate-kit"
    "log"
)

func main() {
    infuraProjectId := "2a1306f1d12f4c109a4d4fb9be46b02e"
    infuraProjectSecret := "fc479a9290b64a84a15fa6544a130218"

    feeRateKit, err := feeratekit.NewFeeRateKit(".", infuraProjectId, infuraProjectSecret)

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
