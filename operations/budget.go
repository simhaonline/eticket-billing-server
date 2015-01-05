package operations

import (
    _ "github.com/lib/pq"
    "database/sql"
)

type Budget struct {
    Merchant string
    Amount int64
}

func (b *Budget) Calculate() (int64, error) {
    conn := NewConnection()

    var amount sql.NullInt64

    ok := conn.QueryRow("select sum(amount) from operations where merchant_id = $1", b.Merchant).Scan(&amount)
    if ok != nil { panic(ok) }

    if amount.Valid {
        b.Amount = amount.Int64
    } else {
        b.Amount = 0
    }

    return b.Amount, nil
}
