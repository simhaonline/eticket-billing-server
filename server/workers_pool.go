package server

import (
	"eticket-billing-server/config"
	glog "github.com/golang/glog"
	"sync"
)

type WorkersPool struct {
	pool        []*Worker
	config      *config.Config
	middlewares MiddlewareChain
	performersMapping PerformerFnMapping
}

var mutex = sync.Mutex{}

func NewWorkersPool(config *config.Config, middlewares MiddlewareChain, mapping PerformerFnMapping) WorkersPool {
	return WorkersPool{config: config, middlewares: middlewares, mapping: mapping}
}

func (wp *WorkersPool) GetWorkerForMerchant(merchant string) *Worker {
	pos := wp.workerPosition(merchant)

	if -1 == pos {
		worker := newWorker(merchant, wp.middlewares, wp.config, wp.performersMapping)
		wp.appendWorker(worker)

		pos = wp.workerPosition(merchant)
		go worker.Serve()
	}

	return wp.pool[pos]
}

func (wp *WorkersPool) appendWorker(worker *Worker) {
	mutex.Lock()
	wp.pool = append(wp.pool, worker)
	mutex.Unlock()
}

func (wp *WorkersPool) StopAll() {
	mutex.Lock()
	// TODO iterate from last and remove it from array
	glog.V(2).Infof("Found %v workers", len(wp.pool))
	for i := 0; i < len(wp.pool); i++ {
		worker := wp.pool[i]
		glog.V(2).Infof("Stopping Worker[%v]...", worker.merchant)
		worker.Stop()
	}
	mutex.Unlock()
}

func (wp WorkersPool) workerPosition(merchant string) int {
	pos := -1
	mutex.Lock()
	for ind, elem := range wp.pool {
		if elem.merchant == merchant {
			pos = ind
		}
	}
	mutex.Unlock()
	return pos
}
