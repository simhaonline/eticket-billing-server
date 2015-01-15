package operations

import (
    pq "github.com/lib/pq"
    "encoding/xml"
    "fmt"
)

type TransactionWithoutCheck struct {
    Transaction
}

func NewTransactionWithoutCheck(data string) *TransactionWithoutCheck {
    r := TransactionWithoutCheck{}
    err := xml.Unmarshal([]byte(data), &r)

    if err != nil {
        fmt.Printf("error: %v", err)
        // TODO pass error
        return nil
    }

    r.OriginXml = data

    return &r
}

func (r *TransactionWithoutCheck) Save() (uint64, error) {
    var id uint64

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

func (r TransactionWithoutCheck) XmlResponse() string {
    tmp := struct {
        TransactionWithoutCheck
        XMLName xml.Name `xml:"answer"`
    }{TransactionWithoutCheck: r}

    tmp.OperationType = "transaction"
    tmp.OriginXml = ""
    output, _ := xml.Marshal(tmp)
    output = append(output, '\n')
    return string(output)
}

func (r TransactionWithoutCheck) ErrorXmlResponse(err error) string {
    error := err.(*TransactionError)

    tmp := struct {
        XMLName xml.Name `xml:"answer"`
        ErrorMessage string `xml:"error>message"`
        ErrorCode string `xml:"error>code"`
        TransactionWithoutCheck
    }{TransactionWithoutCheck: r}

    tmp.ErrorMessage = error.Message
    tmp.ErrorCode = error.Code

    tmp.OperationType = "transaction"
    tmp.OriginXml = ""
    output, _ := xml.Marshal(tmp)
    output = append(output, '\n')
    return string(output)
}
