package feeratekit

import (
    "log"
    "time"
)

type syncer struct {
    storage            *storage
    hsIpfsProvider     ipfsProvider
    globalIpfsProvider ipfsProvider
}

func newSyncer(storage *storage) *syncer {
    return &syncer{
        storage: storage,
        hsIpfsProvider: ipfsProvider{
            baseUrl: hsIpfsBaseUrl,
            timeout: time.Duration(hsIpfsTimeout * time.Second),
        },
        globalIpfsProvider: ipfsProvider{
            baseUrl: globalIpfsBaseUrl,
            timeout: time.Duration(globalIpfsTimeout * time.Second),
        },
    }
}

func (syncer *syncer) syncRates() error {
    rates, err := syncer.hsIpfsProvider.getRates()

    if err == nil {
        return syncer.storage.saveRates(rates)
    }

    log.Println("Could not fetch from HS ipfs:", err)

    rates, err = syncer.globalIpfsProvider.getRates()

    if err == nil {
        return syncer.storage.saveRates(rates)
    }

    log.Println("Could not fetch from global ipfs:", err)

    return err
}
