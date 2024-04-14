package utility

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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
		log.Error("GenerateUUID error #0: ", err)
		return "", err
	}
	return id.String(), nil
}

// GenerateToken генерирует токен
func GenerateToken(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Error("GenerateToken error #0: ", err)
		return "", err
	}

	token := hex.EncodeToString(randomBytes)
	return token, nil
}

// GenerateTokenBase64 генерирует токен в base64
func GenerateTokenBase64(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Error("GenerateToken error #0: ", err)
		return "", err
	}

	// Кодирование случайной строки в base64
	token := base64.URLEncoding.EncodeToString(randomBytes)
	return token, nil
}
