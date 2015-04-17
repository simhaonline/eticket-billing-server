package server

// TODO move to separate package, to make it possible easily to include in different places
// TODO use github.com/jinzhu/gorm
// TODO use gopkg.in/validator.v2
import (
	cfg "eticket-billing-server/config"
	gorm "github.com/jinzhu/gorm"
	pq "gopkg.in/lib/pq.v0"
	"fmt"
)

var currentConfig *cfg.Config

type DbConnection struct {
	Db *gorm.DB
}

func NewConnection(c *cfg.Config) DbConnection {
	connectionString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
		currentConfig.DatabaseUser,
		currentConfig.DatabasePassword,
		currentConfig.DatabaseName)

	db, ok := gorm.Open("postgres", connectionString)
	if ok != nil {
		panic(ok)
	}
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return DbConnection{ Db: &db }
}

func NormalizeDbError(pqError error) OperationError {
	e := pqError.(*pq.Error)
	if e.Code.Name() == "unique_violation" {
		return OperationError{Code: e.Code.Name(), Message: "Duplicate key value violates unique constraint"}
	} else {
		return OperationError{Code: e.Code.Name(), Message: e.Message}
	}

}
