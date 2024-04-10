package main

import (
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/service"
	"cicd-service-go/utility"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// set utc zone for time pkg (as server)
	os.Setenv("TZ", "UTC")

	utility.InitConfig()      // инициалищация конфигурации
	utility.ConfigureLogger() // инициализация логирования
	db.InitInstanceETCD()     // инициализация etcd
	manager.InitManager()
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

	go manager.RunManager() // менеджер управления кластером
	//go schedule.RunCron()		// запуск планировщика задач

	service.Start() // запуск http сервера
}
