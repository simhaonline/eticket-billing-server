package server

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
)

func TestnewWorker(t *testing.T) {
	assert := assert.New(t)

	worker := newWorker(1, "/tmp")
	assert.Equal("*server.Worker", reflect.TypeOf(worker).String(), "NewWorker should return Worker data")
	assert.Equal(worker.merchant, 1, "Merchants doesnt' match")
}
