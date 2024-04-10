package utility

import (
	"github.com/google/uuid"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("CONFIGPATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("InitConfig error #1: fatal error config file: ", err)
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

func GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error("GenerateUUID error #1: ", err)
		return "", err
	}
	return id.String(), nil
}
