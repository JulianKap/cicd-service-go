package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Start() {
	echoInstance := runEchoServer()

	runProfilingServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	//schedule.CloseCronChan <- true

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//taskpkg.CloseAllTasks(ctx)

	if err := echoInstance.Shutdown(ctx); err != nil {
		log.Fatal("Start fatal #0: error shutdown http server for Echo: ", err)
	}
}

func runEchoServer() *echo.Echo {
	echoInstance := startFramework()

	go func() {
		address := fmt.Sprintf("%s:%s", viper.GetString("listen.server.host"), viper.GetString("listen.server.port"))

		if err := echoInstance.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatal("runEchoServer fatal #0: error start http server for Echo: ", err)
		}
	}()

	return echoInstance
}

func runProfilingServer() {
	go func() {
		address := fmt.Sprintf("%s:%s", viper.GetString("listen.pprof.host"), viper.GetString("listen.pprof.port"))

		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatal("runProfilingServer fatal #0: error start http server for profiling: ", err)
		}
	}()
}
