package main

import (
	"cicd-service-go/init/db"
	"cicd-service-go/service"
	"cicd-service-go/utility"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// set utc zone for time pkg (as server)
	os.Setenv("TZ", "UTC")

	utility.InitConfig()
	utility.ConfigureLogger()
	db.InitInstanceETCD()
}

func main() {
	// TODO: добавить трейс в логи
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()
	//go schedule.RunCron()

	service.Start()
}
