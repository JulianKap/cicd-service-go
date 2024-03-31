package utility

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConfigureLogger() {
	logLevel, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		log.Fatalf("Cannot parse log level: %s. Available levels: %s", viper.GetString("logging.level"), log.AllLevels)
	}

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "02.01.2006 - 15:04:05.000",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(logLevel)
	log.Infof("Set %s level logs", logLevel.String())
}
