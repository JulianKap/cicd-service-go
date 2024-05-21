package main

import (
	"cicd-service-go/init/db"
	"cicd-service-go/init/secrets"
	"cicd-service-go/manager"
	"cicd-service-go/schedule"
	"cicd-service-go/service"
	"cicd-service-go/utility"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// set utc zone for time pkg (as server)
	if err := os.Setenv("TZ", "UTC"); err != nil {
		return
	}
	utility.InitConfig()
	utility.ConfigureLogger()
	db.InitInstanceETCD()
	secrets.InitInstanceVault()
	manager.InitManager()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	go manager.RunManager() // запуск менеджера управления кластером
	go schedule.RunCron()   // запуск планировщика задач
	service.Start()         // запуск http сервера
}
