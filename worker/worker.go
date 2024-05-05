package worker

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/scripts"
	"cicd-service-go/taskpkg"
	log "github.com/sirupsen/logrus"
	"sync"
)

// RunWorkerTask запуск выполнения таски
func RunWorkerTask(p pipeline.Pipeline, t taskpkg.Task) (err error) {
	if len(p.Steps) == 0 {
		log.Info("RunWorkerTask #0: null count steps in pipeline (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
		return nil
	}
	var wg sync.WaitGroup

	for _, s := range p.Steps {
		log.Info("RunWorkerTask #1: run step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

		if s.Image == "" {
			log.Warn("RunWorkerTask #2: not image in step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
			//continue
			//var err error
			//return error("")
		}

		// todo: проверить наличие образа, чтобы не пулить
		//wg.Add(1)
		if err := scripts.PullImage(&wg, s.Image); err != nil {
			log.Error("RunWorkerTask #3: ", err)
			return err
		}
		//go scripts.PullImage(&wg, s.Image)
		//wg.Wait()

		//wg.Add(1)
		if err := scripts.RunCommandImage(&wg, s); err != nil {
			log.Error("RunWorkerTask #4: ", err)
			//continue
			return err
		}
		//go scripts.RunCommandImage(&wg, s)
		//wg.Wait()
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
