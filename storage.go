package feeratekit

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "path/filepath"
    "time"
)

type storage struct {
    database *sql.DB
}

func newStorage(dataDir string) *storage {
    dbPath := filepath.Join(dataDir, "fee_rate_kit.db")

    database, _ := sql.Open("sqlite3", dbPath)

    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS fee_rates (coin TEXT PRIMARY KEY ON CONFLICT REPLACE, low INTEGER, medium INTEGER, high INTEGER, timestamp INTEGER)")
    _, _ = statement.Exec()

    return &storage{
        database: database,
    }
}

func (storage storage) save(rate FeeRate) {
    statement, _ := storage.database.Prepare("INSERT INTO fee_rates (coin, low, medium, high, timestamp) VALUES (?, ?, ?, ?, ?)")
    _, _ = statement.Exec(rate.coin, rate.lowPriority, rate.mediumPriority, rate.highPriority, rate.timestamp)
}

func (storage storage) feeRate(coin coin) *FeeRate {
    rows, _ := storage.database.Query("SELECT low, medium, high, timestamp FROM fee_rates WHERE coin = ?", coin)

    var low int64
    var medium int64
    var high int64
    var timestamp int64

    for rows.Next() {
        _ = rows.Scan(&low, &medium, &high, &timestamp)

        //fmt.Printf("DB: low: %v; medium: %v; high: %v; timestamp: %v\n", low, medium, high, timestamp)

        return &FeeRate{
            coin:           coin,
            lowPriority:    low,
            mediumPriority: medium,
            highPriority:   high,
            timestamp:      time.Now().Unix(),
        }
    }

    return nil
}
