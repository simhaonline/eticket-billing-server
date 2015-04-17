package server

type PerformerFn func(*Request) *Request
type PerformerFnMapping map[string]PerformerFn

var mapping PerformerFnMapping

func SetupMapping(mp PerformerFnMapping) {
	mapping = mp
}

func GetMapping(operationName string) PerformerFn {
	return mapping[operationName]
}
