package feeratekit

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "path/filepath"
)

type storage struct {
    database *sql.DB
}

func newStorage(dataDir string) (*storage, error) {
    dbPath := filepath.Join(dataDir, databaseName+".db")

    database, err := sql.Open("sqlite3", dbPath)

    if err != nil {
        return nil, err
    }

    storage := storage{
        database: database,
    }

    err = storage.prepare()

    if err != nil {
        return nil, err
    }

    return &storage, nil
}

func (storage *storage) prepare() error {
    _, err := storage.database.Exec("CREATE TABLE IF NOT EXISTS fee_rates (coin TEXT PRIMARY KEY ON CONFLICT REPLACE, low INTEGER, medium INTEGER, high INTEGER, timestamp INTEGER)")

    return err
}

func (storage *storage) saveRates(rates []*FeeRate) error {
    for _, rate := range rates {
        err := storage.saveRate(rate)

        if err != nil {
            return err
        }
    }

    return nil
}

func (storage *storage) saveRate(rate *FeeRate) error {
    _, err := storage.database.Exec("INSERT INTO fee_rates (coin, low, medium, high, timestamp) VALUES (?, ?, ?, ?, ?)",
        rate.coin, rate.lowPriority, rate.mediumPriority, rate.highPriority, rate.timestamp)
    return err
}

func (storage *storage) feeRate(coin string) *FeeRate {
    row := storage.database.QueryRow("SELECT low, medium, high, timestamp FROM fee_rates WHERE coin = ?", coin)

    rate := FeeRate{coin: coin}

    err := row.Scan(&rate.lowPriority, &rate.mediumPriority, &rate.highPriority, &rate.timestamp)

    if err != nil {
        return nil
    }

    return &rate
}
