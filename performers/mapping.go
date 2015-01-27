package performers

import (
    "eticket-billing-server/request"
)

type PerformerFn func(*request.Request) *request.Request
type PerformerFnMapping map[string]PerformerFn

var mapping PerformerFnMapping

func SetupMapping(mp PerformerFnMapping) {
    mapping = mp
}

func GetMapping(operationName string) PerformerFn {
    return mapping[operationName]
}
