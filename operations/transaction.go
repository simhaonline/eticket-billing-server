package operations

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

type Transaction struct {
	ID              uint     `xml:"-" gorm:"primary_key"`
	XMLName         xml.Name `xml:"request" sql:"-"`
	ApplicationName string   `xml:"application_name"`
	OperationType   string   `xml:"type,attr" sql:"-"`
	Merchant        string   `xml:"merchant" gorm:"column:merchant_id"`
	OperationName   string   `xml:"operation_name"`
	OperationIdent  string   `xml:"operation_ident"`
	Description     string   `xml:"description"`
	Amount          int64    `xml:"amount"`

	OperationCreatedAt customTime `xml:"operation_created_at"`
	OriginXml          string     `xml:",omitempty" sql:"-"`

	Errors []OperationError `sql:"-" xml:"-"`
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

func (t Transaction) TableName() string {
	return "operations"
}

func (r *Transaction) IsPossible() bool {
	budget := Budget{Merchant: r.Merchant}
	amount, _ := budget.Calculate()
	if r.Amount > 0 {
		return true
	} else {
		return (amount+r.Amount > 0)
	}
}

func (r *Transaction) BeforeCreate() (err error) {
	if !r.IsPossible() {
		erro := OperationError{Code: "not_enough_money", Message: "Not enough money for operation"}
		r.Errors = append(r.Errors, erro)
		err = errors.New("Not enough money for operation")
	}
	return
}

func (r *Transaction) XmlResponse() string {
	tmp := struct {
		Transaction
		XMLName xml.Name `xml:"answer"`
	}{Transaction: *r}

	tmp.OperationType = "transaction"
	tmp.OriginXml = ""
	output, _ := xml.Marshal(tmp)
	output = append(output, '\n')
	return string(output)
}

func (r *Transaction) ErrorXmlResponse(err error) string {
	tmp := struct {
		XMLName      xml.Name `xml:"answer"`
		ErrorMessage string   `xml:"error>message"`
		ErrorCode    string   `xml:"error>code"`
		Transaction
	}{Transaction: *r}

	if t := reflect.TypeOf(err); t.String() == "*pq.Error" {
		error := NormalizeDbError(err)
		tmp.ErrorMessage = error.Message
		tmp.ErrorCode = error.Code
	} else if len(r.Errors) > 0 {
		error := r.Errors[0]
		tmp.ErrorMessage = error.Message
		tmp.ErrorCode = error.Code
	}

	tmp.OperationType = "transaction"
	tmp.OriginXml = ""
	output, _ := xml.Marshal(tmp)
	output = append(output, '\n')
	return string(output)
}
