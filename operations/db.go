package operations

import (
	"database/sql"
	cfg "eticket-billing-server/config"
	"fmt"
	_ "gopkg.in/lib/pq.v0"
)

var currentConfig *cfg.Config

func SetupConnections(c *cfg.Config) {
	currentConfig = c
}

func NewConnection() *sql.DB {
	connectionString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
		currentConfig.DatabaseUser,
		currentConfig.DatabasePassword,
		currentConfig.DatabaseName)

	db, ok := sql.Open("postgres", connectionString)
	if ok != nil {
		panic(ok)
	}

	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
