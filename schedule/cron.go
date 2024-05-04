package schedule

import (
	"cicd-service-go/manager"
	"cicd-service-go/worker/hub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	CloseCronChan chan bool
)

func init() {
	CloseCronChan = make(chan bool)
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

// сделать исключение для полей в структурах, чтобы можно было доавбить больше полей
// все таки сделать упрощение для воркера и просто использовать chan
//

func runSchedule() error {
	workerHub := hub.GetHub()

	workerHub.ChanLock()
	defer workerHub.ChanUnlock()

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

	//workerHub := hub.GetHub()
	//
	//workerHub.ChanLock()
	//defer workerHub.ChanUnlock()

	if err := runTasksWorker(); err != nil {
		log.Error("runSchedule #1: ", err)
		return err
	}

	// Получаем свои задачи
	// Запускаем
	// Использовать алгоритм пайплайн
	//

	return nil
}
