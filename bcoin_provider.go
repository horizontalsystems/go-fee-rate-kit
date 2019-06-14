package feeratekit

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "time"
)

type bcoinProvider struct {
    baseUrl string
}

func (provider *bcoinProvider) getBitcoinRate() (*FeeRate, error) {
    low, err := provider.getRateValue(15)
    medium, err := provider.getRateValue(6)
    high, err := provider.getRateValue(1)

    if err != nil {
        return nil, err
    }

    rate := FeeRate{
        coin:           bitcoin,
        lowPriority:    low,
        mediumPriority: medium,
        highPriority:   high,
        timestamp:      time.Now().Unix(),
    }

    return &rate, err
}

func (provider *bcoinProvider) getRateValue(numberOfBlocks int) (int64, error) {
    requestBody, err := json.Marshal(map[string]interface{}{
        "method": "estimatesmartfee",
        "params": []int{numberOfBlocks},
    })

    if err != nil {
        return -1, err
    }

    response, err := http.Post(provider.baseUrl, "application/json", bytes.NewBuffer(requestBody))

    if err != nil {
        return -1, err
    }

    defer response.Body.Close()

    resultMap := make(map[string]interface{})

    err = json.NewDecoder(response.Body).Decode(&resultMap)

    if err != nil {
        return -1, err
    }

    result, success := resultMap["result"].(map[string]interface{})

    if !success {
        return -1, errors.New("type conversion error: 'result' field")
    }

    fee, success := result["fee"].(float64)

    if !success {
        return -1, errors.New("type conversion error: 'result'['fee'] field")
    }

    if fee < 0 {
        return -1, errors.New("invalid fee: < 0")
    }

    convertedFee := int64(fee * 100000000 / 1024)

    return convertedFee, nil
}
