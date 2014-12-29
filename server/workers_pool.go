package server

import (
	"sync"
)

type WorkersPool []*Worker

var workersPoolInstance WorkersPool

var mutex = sync.Mutex{}

func NewWorkersPool() *WorkersPool {
	if workersPoolInstance == nil {
		workersPoolInstance = make(WorkersPool, 0)
	}

	return &workersPoolInstance
}

func (wp *WorkersPool) GetWorkerForMerchant(merchant int) *Worker {
	pos := wp.workerPosition(merchant)

	if -1 == pos {
		worker := newWorker(merchant, "/tmp")
		wp.appendWorker(worker)

		pos = wp.workerPosition(merchant)
		go worker.Serve()
	}

	return (*wp)[pos]
}

func (wp *WorkersPool) appendWorker(worker *Worker) {
	mutex.Lock()
	workersPoolInstance = append(*wp, worker)
	mutex.Unlock()
}

func (wp WorkersPool) workerPosition(merchant int) int {
	pos := -1
	mutex.Lock()
	for ind, elem := range wp {
		if elem.merchant == merchant {
			pos = ind
		}
	}
	mutex.Unlock()
	return pos
}
