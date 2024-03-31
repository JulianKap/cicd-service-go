package main

import (
	"cicd-service-go/service"
	"cicd-service-go/utility"
	"os"
)

func init() {
	// set utc zone for time pkg (as server)
	os.Setenv("TZ", "UTC")

	utility.InitConfig()
	utility.ConfigureLogger()
	//etcd.InitInstanceDB()
}

func main() {
	//go schedule.RunCron()

	service.Start()
}
