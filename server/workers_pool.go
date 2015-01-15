package server

import (
    "sync"
    glog "github.com/golang/glog"
    "eticket-billing-server/config"
)

type WorkersPool []*Worker

var workersPoolInstance WorkersPool

var mutex = sync.Mutex{}

func GetWorkersPool() *WorkersPool {
    if workersPoolInstance == nil {
        workersPoolInstance = make(WorkersPool, 0)
    }

    return &workersPoolInstance
}

func (wp *WorkersPool) GetWorkerForMerchant(merchant string) *Worker {
    pos := wp.workerPosition(merchant)

    if -1 == pos {
        config := config.GetConfig()
        worker := newWorker(merchant, config.RequestLogDir)
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

func (wp *WorkersPool) StopAll() {
    mutex.Lock()
    // TODO iterate from last and remove it from array
    glog.V(2).Infof("Found %v workers", len(workersPoolInstance))
    for i := 0; i < len(workersPoolInstance); i++ {
        worker := (*wp)[i]
        glog.V(2).Infof("Stopping Worker[%v]...", worker.merchant)
        worker.Stop()
    }
    mutex.Unlock()
}

func (wp WorkersPool) workerPosition(merchant string) int {
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
