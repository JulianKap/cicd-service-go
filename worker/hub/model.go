package hub

import "cicd-service-go/worker"

type workerHub struct {
	runningWorker worker.IWorker
	controlChan   chan int `json:"-"`
}
