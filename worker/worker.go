package worker

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/scripts"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// RunWorkerTask запуск выполнения таски
func RunWorkerTask(j sources.Job, p pipeline.Pipeline, t taskpkg.Task) ([]taskpkg.SubTaskResult, error) {
	if len(p.Steps) == 0 {
		log.Info("RunWorkerTask #0: null count steps in pipeline (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
		return nil, nil
	}
	var wg sync.WaitGroup

	subTask := taskpkg.PrepareSubTasksStruct(p)

	for i, s := range p.Steps {
		log.Info("Step: ", s.Name, " START (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

		if s.Image == "" {
			log.Warn("RunWorkerTask #2: not image in step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
			// todo: сделать вариант запуска локально, если нет image

			subTask[i].Status = taskpkg.StepFailed
			subTask[i].Message = "Docker image is missing"

			return subTask, errors.New("image is null")
		} else {
			if s.Branch != j.Branch {
				log.Info("RunWorkerTask #3: null count steps in pipeline (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

				subTask[i].Status = taskpkg.StepSkipped
				subTask[i].Message = "The repository branch is not suitable"
				continue
			}

			commands, err := taskpkg.PrepareStepCommand(j, s)
			if err != nil {
				log.Error("RunWorkerTask #4: ", err)

				subTask[i].Status = taskpkg.StepFailed
				subTask[i].Message = "Error preparing commands to run"
				return subTask, err
			}

			if commands == "" {
				log.Warn("RunWorkerTask #5: command is null for step=", s.Name, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")

				subTask[i].Status = taskpkg.StepSkipped
				subTask[i].Message = "Missing commands to run in step"
				continue
			}

			ok, err := scripts.CheckImageExists(s.Image)
			if err != nil {
				log.Error("RunWorkerTask #6: ", err)
				//return err
			}

			if !ok {
				// todo: проверить наличие образа, чтобы не пулить
				//wg.Add(1)
				if err := scripts.PullImage(&wg, s.Image); err != nil {
					log.Error("RunWorkerTask #7: ", err)

					subTask[i].Status = taskpkg.StepFailed
					subTask[i].Message = "Error getting docker image: " + s.Image
					return subTask, err
				}
				//wg.Wait()
			}

			start := time.Now()
			subTask[i].Status = taskpkg.StepRunning

			//wg.Add(1)
			if err := scripts.RunCommandImage(&wg, s.Image, commands); err != nil {
				log.Error("RunWorkerTask #8: ", err)

				subTask[i].Status = taskpkg.StepFailed
				subTask[i].Message = "Errors when executing step commands"
				return subTask, err
			}

			elapsed := time.Since(start)
			subTask[i].ElapsedSec = int(elapsed.Seconds())
			subTask[i].Status = taskpkg.StepComplited

			//wg.Wait()
		}

		log.Info("Step: ", s.Name, " DONE (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
	}

	return subTask, nil
}
