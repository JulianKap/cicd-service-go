package worker

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/scripts"
	"cicd-service-go/taskpkg"
	log "github.com/sirupsen/logrus"
	"sync"
)

// RunTask запуск выполнения таски
func RunTask(p pipeline.Pipeline, t taskpkg.Task) (err error) {
	if len(p.Steps) == 0 {
		log.Info("RunTask #0: null count steps in pipeline project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID)
		return nil
	}

	for _, s := range p.Steps {
		log.Info("RunTask #1: run step=", s.Name, " project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID)

		if s.Image == "" {
			log.Warn("RunTask #2: not image in step=", s.Name, " project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID)
			continue
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go scripts.PullImage(&wg, s.Image)
		wg.Wait()

		wg.Add(1)
		go scripts.RunCommandImage(&wg, s)
		wg.Wait()
	}

	// todo: Использовать патерн пайплайн

	return nil
}

//
//import "cicd-service-go/taskpkg"
//
//type IWorker interface {
//	CheckTaskName(taskName string) error
//	StartWork() (task taskpkg.ITask, isRunning bool, err error)
//	StatusWork() (task taskpkg.ITask, isRunning bool, err error)
//	SetOpts(scheduleID int, taskName string) error
//}
//
//type IWorkerHub interface {
//	GetWorker(workerName string) IWorker
//	GetStatusRunningWork() (taskpkg.ITask, bool, error)
//	GetRunningWorker() IWorker
//	SetRunningWorker(IWorker)
//	DelRunningWork()
//	ChanLock()
//	ChanUnlock()
//}
