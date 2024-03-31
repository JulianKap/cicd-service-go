package utility

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("CONFIGPATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Fatal error config file: ", err)
	}
}

func StringPtr(str string) *string {
	if str == "" {
		return nil
	}

	return &str
}

func StringPtrToString(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}
