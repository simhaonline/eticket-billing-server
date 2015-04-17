package server

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

type TransactionWithoutCheck struct {
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

	Db *DbConnection `xml:"-"`
}

func NewTransactionWithoutCheck(data string, db *DbConnection) *TransactionWithoutCheck {
	r := TransactionWithoutCheck{Db: db}
	err := xml.Unmarshal([]byte(data), &r)

	if err != nil {
		fmt.Printf("error: %v", err)
		// TODO pass error
		return nil
	}

	r.OriginXml = data

	return &r
}

func (r *TransactionWithoutCheck) XmlResponse() string {
	tmp := struct {
		TransactionWithoutCheck
		XMLName xml.Name `xml:"answer"`
	}{TransactionWithoutCheck: *r}

	tmp.OperationType = "transaction"
	tmp.OriginXml = ""
	output, _ := xml.Marshal(tmp)
	output = append(output, '\n')
	return string(output)
}

func (r *TransactionWithoutCheck) ErrorXmlResponse(err error) string {
	// TODO remove duplication
	tmp := struct {
		XMLName      xml.Name `xml:"answer"`
		ErrorMessage string   `xml:"error>message"`
		ErrorCode    string   `xml:"error>code"`
		TransactionWithoutCheck
	}{TransactionWithoutCheck: *r}

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
