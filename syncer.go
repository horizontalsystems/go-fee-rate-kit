package feeratekit

import (
    "log"
    "time"
)

type syncer struct {
    delegate           syncerDelegate
    storage            *storage
    hsIpfsProvider     *ipfsProvider
    globalIpfsProvider *ipfsProvider
    bcoinProvider      *bcoinProvider
    infuraProvider     *infuraProvider
}

type syncerDelegate interface {
    didSyncRates()
}

func newSyncer(storage *storage, infuraProjectId string, infuraProjectSecret string) *syncer {
    var infura *infuraProvider

    if infuraProjectId != "" {
        infura = &infuraProvider{
            projectId:     infuraProjectId,
            projectSecret: infuraProjectSecret,
        }
    }

    return &syncer{
        storage: storage,
        hsIpfsProvider: &ipfsProvider{
            baseUrl: hsIpfsBaseUrl,
            timeout: time.Duration(hsIpfsTimeout * time.Second),
        },
        globalIpfsProvider: &ipfsProvider{
            baseUrl: globalIpfsBaseUrl,
            timeout: time.Duration(globalIpfsTimeout * time.Second),
        },
        bcoinProvider: &bcoinProvider{
            baseUrl: bcoinBaseUrl,
        },
        infuraProvider: infura,
    }
}

func (syncer *syncer) syncRates() {
    go syncer.syncFromHsIpfs()
}

func (syncer *syncer) syncFromHsIpfs() {
    rates, err := syncer.hsIpfsProvider.getRates()

    if err != nil {
        log.Println("HS IPFS failed:", err)

        go syncer.syncFromGlobalIpfs()

        return
    }

    log.Println("HS IPFS success")

    syncer.handleRates(rates)
}

func (syncer *syncer) syncFromGlobalIpfs() {
    rates, err := syncer.globalIpfsProvider.getRates()

    if err != nil {
        log.Println("Global IPFS failed:", err)

        go syncer.syncFromBcoin()
        go syncer.syncFromInfura()

        return
    }

    log.Println("Global IPFS success")

    syncer.handleRates(rates)
}

func (syncer *syncer) syncFromBcoin() {
    rate, err := syncer.bcoinProvider.getBitcoinRate()

    if err != nil {
        log.Println("Bcoin failed:", err)
        return
    }

    log.Println("Bcoin success")

    syncer.handleRates([]*FeeRate{rate})
}

func (syncer *syncer) syncFromInfura() {
    if syncer.infuraProvider == nil {
        return
    }

    rate, err := syncer.infuraProvider.getEthereumRate()

    if err != nil {
        log.Println("Infura failed:", err)
        return
    }

    log.Println("Infura success")

    syncer.handleRates([]*FeeRate{rate})
}

func (syncer *syncer) handleRates(rates []*FeeRate) {
    err := syncer.storage.saveRates(rates)

    if err != nil {
        log.Println("Could not save rates:", err)
        return
    }

    syncer.delegate.didSyncRates()
}
