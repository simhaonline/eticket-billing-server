package operations

import (
    _ "github.com/lib/pq"
    "database/sql"
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
    // TODO move connection's params to config file
    db, ok := sql.Open("postgres", "user=ignar dbname=eticket_billing sslmode=disable")
    if ok != nil { panic (ok) }

    err := db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return db
}
