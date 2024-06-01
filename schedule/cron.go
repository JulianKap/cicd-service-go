package schedule

import (
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/pipeline"
	"cicd-service-go/taskpkg"
	"cicd-service-go/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	CloseCronChan chan bool
	controlChan   chan int
)

func init() {
	CloseCronChan = make(chan bool)
}

func сhanLock() {
	controlChan <- 1
}

func сhanUnlock() {
	<-controlChan
}

func RunCron() {
	countParallelTasks := viper.GetInt("schedule.tasks_parallel_workers")
	if countParallelTasks < 1 {
		countParallelTasks = 5 // default
	}

	if countParallelTasks > 100 {
		countParallelTasks = 100
	}

	controlChan = make(chan int, countParallelTasks)

	cronTimeSeconds := viper.GetInt("schedule.cron_timer_s")
	if cronTimeSeconds < 1 {
		cronTimeSeconds = 60
	}

	scheduleTicker := time.NewTicker(time.Duration(cronTimeSeconds) * time.Second)

	for {
		select {
		case <-CloseCronChan:
			log.Info("RunCron #0: exit RunCron")
			return
		case <-scheduleTicker.C:
			if err := runSchedule(); err != nil {
				log.Errorln("RunCron #1: ", err)
			}
		case <-scheduleTicker.C:
			if err := runScheduleWorker(); err != nil {
				log.Errorln("RunCron #2: ", err)
			}
		}
	}
}

func runSchedule() error {
	if manager.MemberInfo.ReadOnly {
		log.Debug("runSchedule #0: service in read only mode")
		return nil
	}

	if manager.MemberInfo.Role == manager.MasterRole {
		log.Debug("runSchedule #1: run scheduler as MASTER")

		standalone, err := tasksScheduler()
		if err != nil {
			log.Error("runSchedule #2: ", err)
			return err
		}

		if standalone {
			log.Debug("runSchedule #3: continuation of work in mode STANDALONE")
			manager.MemberInfo.Standalone = true
		} else {
			manager.MemberInfo.Standalone = false
		}
	}
	return nil
}

func runScheduleWorker() error {
	if manager.MemberInfo.ReadOnly {
		log.Debug("runScheduleWorker #0: service in read only mode")
		return nil
	}

	if manager.MemberInfo.Role == manager.MasterRole && !manager.MemberInfo.Standalone {
		log.Debug("runScheduleWorker #1: in master rule")
		return nil
	}

	go func() {
		сhanLock()
		defer сhanUnlock()

		log.Debug("runScheduleWorker #2: run scheduler for Worker")

		tasks, err := tasksSchedulerWorker()
		if err != nil {
			log.Error("runScheduleWorker #3: ", err)
			//return err
			return
		}

		if len(tasks.Tasks) == 0 {
			log.Debug("runScheduleWorker #4: not found actual tasks for ", manager.MemberInfo.UUID)
			//return nil
			return
		}

		// Запуск выполнения задач
		for _, t := range tasks.Tasks {
			j, err := getJobEtcd(*t)
			if err != nil {
				log.Error("runScheduleWorker #5: ", err)
				continue
			}

			// Получение пайплайна
			p, err := pipeline.GetPipeline(j)
			if err != nil {
				log.Error("runScheduleWorker #6: ", err)
				//todo: возможно тоже нужно отмечать такие задания в etcd с другим статусом
				continue
			}

			start := time.Now()
			t.Status.RunningAt = &start
			t.Status.Status = taskpkg.Running

			// Отмечаем, что задание уже запущено
			if err := updateTaskForWorker(db.InstanceETCD, manager.MemberInfo, t); err != nil {
				log.Error("runScheduleWorker #7: ", err)
				continue
			}

			// Запуск выполнения пайплайна
			subTasks, err := worker.RunWorkerTask(j, p, *t)
			if err != nil {
				log.Error("runScheduleWorker #8: ", err)

				t.Status.Status = taskpkg.Failed
				t.Status.Message = err.Error()
				t.Status.RetryCount++
			} else {
				t.Status.Status = taskpkg.Completed
			}

			elapsed := time.Since(start)
			t.Status.ElapsedSec = int(elapsed.Seconds())
			t.Status.Steps = subTasks

			if err := updateTaskForWorker(db.InstanceETCD, manager.MemberInfo, t); err != nil {
				log.Error("runScheduleWorker #9: ", err)
				continue
			}
		}
	}()

	return nil
}
