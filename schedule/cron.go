package schedule

import (
	"cicd-service-go/manager"
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
	scheduleTicker := time.NewTicker(time.Duration(viper.GetInt("schedule.cron_timer_s")) * time.Second)

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
	if manager.MemberInfo.Master {
		log.Debug("runSchedule #0: run scheduler as MASTER")

	} else {
		log.Debug("runSchedule #1: run scheduler as Worker")

	}

	return nil
}
