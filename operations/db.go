package operations

import (
    _ "github.com/lib/pq"
    "database/sql"
    "eticket-billing-server/config"
    "fmt"
)

func NewConnection() *sql.DB {
    config := config.GetConfig()

    connectionString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DatabaseUser, config.DatabasePassword, config.DatabaseName)

    db, ok := sql.Open("postgres", connectionString)
    if ok != nil { panic (ok) }

    err := db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return db
}
