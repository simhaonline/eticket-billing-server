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

func GetWorkerForMerchant(merchant int) *Worker {
	pos := workerPosition(merchant)

	if -1 == pos {
		worker := newWorker(merchant, "/tmp")
		appendWorker(worker)

		pos = workerPosition(merchant)
		go worker.Serve()
	}

	return workersPoolInstance[pos]
}

func appendWorker(worker *Worker) {
	mutex.Lock()
	workersPoolInstance = append(workersPoolInstance, worker)
	mutex.Unlock()
}

func workerPosition(merchant int) int {
	pos := -1
	mutex.Lock()
	for ind, elem := range workersPoolInstance {
		if elem.merchant == merchant {
			pos = ind
		}
	}
	mutex.Unlock()
	return pos
}
