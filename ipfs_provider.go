package feeratekit

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type ipfsProvider struct {
    baseUrl string
    timeout time.Duration
}

type ipfsResponse struct {
    Time  int64
    Rates map[string]ipfsRateResponse
}

type ipfsRateResponse struct {
    Low    int64 `json:"low_priority"`
    Medium int64 `json:"medium_priority"`
    High   int64 `json:"high_priority"`
}

func (provider *ipfsProvider) getRates() ([]FeeRate, error) {
    ipfsUrl := provider.baseUrl + ipfsPath

    log.Printf("Requesting %v", ipfsUrl)

    client := http.Client{
        Timeout: provider.timeout,
    }

    response, err := client.Get(ipfsUrl)

    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    ipfsResponse := ipfsResponse{}

    err = json.NewDecoder(response.Body).Decode(&ipfsResponse)

    if err != nil {
        return nil, err
    }

    rates := make([]FeeRate, 0)

    for coin, rateResponse := range ipfsResponse.Rates {
        //log.Printf("%v: %+v", coin, ipfsRateResponse)

        rate := FeeRate{
            coin:           coin,
            lowPriority:    rateResponse.Low,
            mediumPriority: rateResponse.Medium,
            highPriority:   rateResponse.High,
            timestamp:      ipfsResponse.Time,
        }

        rates = append(rates, rate)
    }

    return rates, nil
}
