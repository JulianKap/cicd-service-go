package worker

import "cicd-service-go/taskpkg"

type IWorker interface {
	CheckTaskName(taskName string) error
	StartWork() (task taskpkg.ITask, isRunning bool, err error)
	StatusWork() (task taskpkg.ITask, isRunning bool, err error)
	SetOpts(scheduleID int, taskName string) error
}

type IWorkerHub interface {
	GetWorker(workerName string) IWorker
	GetStatusRunningWork() (taskpkg.ITask, bool, error)
	GetRunningWorker() IWorker
	SetRunningWorker(IWorker)
	DelRunningWork()
	ChanLock()
	ChanUnlock()
}
