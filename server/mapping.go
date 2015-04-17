package server

type PerformerFn func(*Request, *DbConnection)*performerType
type PerformerFnMapping map[string]PerformerFn
