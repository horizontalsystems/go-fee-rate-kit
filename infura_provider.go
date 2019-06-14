package feeratekit

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "time"
)

type infuraProvider struct {
    projectId     string
    projectSecret string
}

func (provider *infuraProvider) getEthereumRate() (*FeeRate, error) {
    url := infuraBaseUrl + "/" + provider.projectId

    requestBody, err := json.Marshal(map[string]interface{}{
        "id":      "1",
        "jsonrpc": "2.0",
        "method":  "eth_gasPrice",
        "params":  make([]interface{}, 0),
    })

    if err != nil {
        return nil, err
    }

    request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

    if err != nil {
        return nil, err
    }

    request.Header.Set("Accept", "application/json")

    projectSecret := provider.projectSecret
    if projectSecret != "" {
        request.Header.Set("Authorization", basicAuth("", projectSecret))
    }

    response, err := http.DefaultClient.Do(request)

    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    resultMap := make(map[string]interface{})

    err = json.NewDecoder(response.Body).Decode(&resultMap)

    if err != nil {
        return nil, err
    }

    feeHex, success := resultMap["result"].(string)

    if !success {
        return nil, errors.New("type conversion error: 'result' field")
    }

    fee, err := strconv.ParseInt(feeHex, 0, 64)

    if err != nil {
        return nil, err
    }

    rate := FeeRate{
        coin:           ethereum,
        lowPriority:    fee / 2,
        mediumPriority: fee,
        highPriority:   fee * 2,
        timestamp:      time.Now().Unix(),
    }

    return &rate, nil
}

func basicAuth(username, password string) string {
    auth := username + ":" + password
    return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
