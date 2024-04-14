package schedule

import (
	"cicd-service-go/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var CloseCronChan chan bool

func init() {
	CloseCronChan = make(chan bool)
}

func RunCron() {
	scheduleTicker := time.NewTicker(time.Duration(viper.GetInt("schedule.cron_timer_s")) * time.Second)

	for {
		select {
		case <-CloseCronChan:
			log.Info("RunCron info #0: сlose RunCron")
			return
		case <-scheduleTicker.C:
			if err := runSchedule(); err != nil {
				log.Errorln("RunCron error #1: ", err)
			}
		}
	}
}

func runSchedule() error {
	if manager.MemberInfo.Master {
		log.Debug("runSchedule debug #0: run scheduler as MASTER")

	} else {
		log.Debug("runSchedule debug #1: run scheduler as Worker")

	}

	return nil
}
