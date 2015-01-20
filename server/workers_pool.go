package server

import (
    "sync"
    glog "github.com/golang/glog"
    "eticket-billing-server/config"
)

type WorkersPool struct {
    pool []*Worker
    config *config.Config
}

var mutex = sync.Mutex{}

func NewWorkersPool(config *config.Config) WorkersPool {
    return WorkersPool{config: config}
}

func (wp *WorkersPool) GetWorkerForMerchant(merchant string) *Worker {
    pos := wp.workerPosition(merchant)

    if -1 == pos {
        worker := newWorker(merchant, wp.config.RequestLogDir)
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
