package worker

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/scripts"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

// RunWorkerTask запуск выполнения таски
func RunWorkerTask(j sources.Job, p pipeline.Pipeline, t taskpkg.Task) (err error) {
	if len(p.Steps) == 0 {
		log.Info("RunWorkerTask #0: null count steps in pipeline (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
		return nil
	}
	var wg sync.WaitGroup

	for _, s := range p.Steps {
		log.Info("RunWorkerTask #1: run step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

		if s.Image == "" {
			log.Warn("RunWorkerTask #2: not image in step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

			// todo: сделать вариант запуска локально, если нет image

			return errors.New("image is null")
		} else {
			commands, err := taskpkg.PrepareStepCommand(j, s)
			if err != nil {
				log.Error("RunWorkerTask #3: ", err)
				return err
			}

			if commands == "" {
				log.Warn("RunWorkerTask #4: command is null step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
				return errors.New("commands is null")
			}

			// todo: проверить наличие образа, чтобы не пулить
			//wg.Add(1)
			if err := scripts.PullImage(&wg, s.Image); err != nil {
				log.Error("RunWorkerTask #5: ", err)
				return err
			}
			//go scripts.PullImage(&wg, s.Image)
			//wg.Wait()

			//wg.Add(1)
			if err := scripts.RunCommandImage(&wg, s.Image, commands); err != nil {
				log.Error("RunWorkerTask #6: ", err)
				//continue
				return err
			}
			//go scripts.RunCommandImage(&wg, s)
			//wg.Wait()
		}
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
