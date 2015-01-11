package operations

import (
    _ "github.com/lib/pq"
    "database/sql"
    "eticket-billing/config"
    "fmt"
)

/*
create table operations (
  id serial primary key,
  merchant_id varchar(20),
  operation_ident varchar(20),
  description text,
  amount integer,
  xml_data xml,
  operation_created_at timestamp without time zone,
  created_at timestamp without time zone default now(),
  updated_at timestamp without time zone,
  UNIQUE (merchant_id, operation_ident)
)
*/

func NewConnection() *sql.DB {
    config := config.GetConfig()

    connectionString := fmt.Sprintf("user=%v dbname=%v sslmode=disable", config.DatabaseUser, config.DatabaseName)

    db, ok := sql.Open("postgres", connectionString)
    if ok != nil { panic (ok) }

    err := db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return db
}
