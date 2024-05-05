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
	controlChan = make(chan int, 1)
}

func сhanLock() {
	controlChan <- 1
}

func сhanUnlock() {
	<-controlChan
}

func RunCron() {
	cronTimeSeconds := viper.GetInt("schedule.cron_timer_s")
	if cronTimeSeconds < 1 {
		cronTimeSeconds = 60 // default
	}

	scheduleTicker := time.NewTicker(time.Duration(cronTimeSeconds) * time.Second)

	for {
		select {
		case <-CloseCronChan:
			log.Info("RunCron #0: сlose RunCron")
			return
		case <-scheduleTicker.C:
			if err := runSchedule(); err != nil {
				log.Errorln("RunCron #1: ", err)
			}
		}
	}
}

func runSchedule() error {
	сhanLock()
	defer сhanUnlock()

	if manager.MemberInfo.Master {
		log.Debug("runSchedule #0: run scheduler as MASTER")

		standalone, err := tasksScheduler()
		if err != nil {
			log.Error("runSchedule #1: ", err)
			return err
		}

		if !standalone {
			return nil
		}

		log.Debug("runSchedule #2: continuation of work in mode STANDALONE")
	}

	log.Debug("runSchedule #3: run scheduler for Worker")

	err, tasks := tasksSchedulerWorker()
	if err != nil {
		log.Error("runSchedule #4: ", err)
		return err
	}

	if len(tasks.Tasks) == 0 {
		log.Debug("runSchedule #5: not found actual tasks for ", manager.MemberInfo.UUID)
		return nil
	}

	// Запуск выполнения задач
	for _, t := range tasks.Tasks {
		if t.NumberOfRetriesOnError < t.Status.RetryCount {
			continue
		}

		j, err := getJobEtcd(t)
		if err != nil {
			log.Error("runSchedule #6: ", err)
			continue
		}

		// Получение пайплайна
		p, err := pipeline.GetPipeline(j)
		if err != nil {
			log.Error("runSchedule #7: ", err)
			//todo: возможно тоже нужно отмечать такие задания в etcd с другим статусом
			continue
		}

		tm := time.Now()
		t.CreateAt = &tm

		// Запуск выполнения пайплайна
		if err := worker.RunWorkerTask(p, t); err != nil {
			log.Error("runSchedule #8: ", err)
			t.Status.Status = taskpkg.Failed
			t.Status.Message = err.Error()
			t.Status.RetryCount++

			//continue
		} else {
			t.Status.Status = taskpkg.Completed
		}

		if err := updateTaskForWorker(db.InstanceETCD, manager.MemberInfo, &t); err != nil {
			log.Error("runSchedule #9: ", err)
			continue
		}
	}

	return nil
}
