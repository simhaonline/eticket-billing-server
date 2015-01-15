package operations

import (
    pq "github.com/lib/pq"
    "encoding/xml"
    "fmt"
)

type Transaction struct {
    XMLName xml.Name `xml:"request"`
    OperationType string `xml:"type,attr"`
    Merchant string `xml:"merchant"`
    OperationIdent string `xml:"operation_ident"`
    Description string `xml:"description"`
    Amount int64 `xml:"amount"`

    OperationCreatedAt customTime `xml:"operation_created_at"`
    OriginXml string `xml:",omitempty"`
}

func NewTransaction(data string) *Transaction {
    r := Transaction{}
    err := xml.Unmarshal([]byte(data), &r)

    if err != nil {
        fmt.Printf("error: %v", err)
        // TODO pass error
        return nil
    }

    r.OriginXml = data

    return &r
}

func (r *Transaction) IsPossible() bool {
    budget := Budget{Merchant: r.Merchant}
    amount, _ := budget.Calculate()
    if r.Amount > 0 {
        return true
    } else {
        return (amount + r.Amount > 0)
    }
}

func (r *Transaction) Save() (uint64, error) {
    var id uint64

    if !r.IsPossible() {
        return 0, &TransactionError{Code: "not_enough_money", Message: "Not enough money for operation"}
    } else {
        conn := NewConnection()
        defer conn.Close()
        ok := conn.QueryRow(`INSERT INTO operations (merchant_id, operation_ident, description, amount, operation_created_at, xml_data)
                         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
            r.Merchant, r.OperationIdent, r.Description, r.Amount, r.OperationCreatedAt, r.OriginXml).Scan(&id)

        if err, ok := ok.(*pq.Error); ok {
            if "unique_violation" == err.Code.Name() {
                fmt.Println("Duplication warning")
                return 0, &TransactionError{Code: "duplication", Message: "Duplication of operation"}
            } else {
                panic(err)
            }
        }
        return id, nil
    }
}

func (r Transaction) XmlResponse() string {
    tmp := struct {
        Transaction
        XMLName xml.Name `xml:"answer"`
    }{Transaction: r}

    tmp.OperationType = "transaction"
    tmp.OriginXml = ""
    output, _ := xml.Marshal(tmp)
    output = append(output, '\n')
    return string(output)
}

func (r Transaction) ErrorXmlResponse(err error) string {
    error := err.(*TransactionError)

    tmp := struct {
        XMLName xml.Name `xml:"answer"`
        ErrorMessage string `xml:"error>message"`
        ErrorCode string `xml:"error>code"`
        Transaction
    }{Transaction: r}

    tmp.ErrorMessage = error.Message
    tmp.ErrorCode = error.Code

    tmp.OperationType = "transaction"
    tmp.OriginXml = ""
    output, _ := xml.Marshal(tmp)
    output = append(output, '\n')
    return string(output)
}
