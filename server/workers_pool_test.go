package server

import (
    "testing"
    "reflect"
    "github.com/stretchr/testify/assert"
)

func TestNewWorkersPool(t *testing.T) {
    assert := assert.New(t)

    pool := *NewWorkersPool()
    assert.Equal("server.WorkersPool", reflect.TypeOf(pool).String(), "NewWorkersPool should return pool of workers")
    assert.Equal(0, len(pool), "Pool of workers must be empty")
}

func TestGetWorkerForMerchant(t *testing.T) {
    assert := assert.New(t)

    pool := NewWorkersPool()


    worker := pool.GetWorkerForMerchant("10")
    assert.Equal("10", worker.merchant, "It should create new worker for merchant")
    assert.Equal(1, len(workersPoolInstance), "It should change length of pool by 1")

    worker = pool.GetWorkerForMerchant("10")
    assert.Equal("10", worker.merchant, "It should fetch worker for merchant")
    assert.Equal(1, len(workersPoolInstance), "It should not change length of pool after second call")

    worker = pool.GetWorkerForMerchant("20")
    assert.Equal("20", worker.merchant, "It should fetch worker for merchant")
    assert.Equal(2, len(workersPoolInstance), "It should not change length of pool after second call")
}
