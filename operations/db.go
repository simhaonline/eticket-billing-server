package operations

// TODO move to separate package, to make it possible easily to include in different places
// TODO use github.com/jinzhu/gorm
// TODO use gopkg.in/validator.v2
import (
	cfg "eticket-billing-server/config"
	"fmt"
	gorm "github.com/jinzhu/gorm"
	pq "gopkg.in/lib/pq.v0"
)

var db *gorm.DB

var currentConfig *cfg.Config

func SetupConnections(c *cfg.Config) {
	currentConfig = c
}

// TODO Maybe connect every time?
func NewConnection() *gorm.DB {
	connectionString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
		currentConfig.DatabaseUser,
		currentConfig.DatabasePassword,
		currentConfig.DatabaseName)

	dbd, ok := gorm.Open("postgres", connectionString)
	if ok != nil {
		panic(ok)
	}
	db = &dbd
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}

func GetDB() *gorm.DB {
	return db
}

func NormalizeDbError(pqError error) OperationError {
	e := pqError.(*pq.Error)
	if e.Code.Name() == "unique_violation" {
		return OperationError{Code: e.Code.Name(), Message: "Duplicate key value violates unique constraint"}
	} else {
		return OperationError{Code: e.Code.Name(), Message: e.Message}
	}

}
