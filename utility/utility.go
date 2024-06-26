package utility

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("CONFIGPATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("InitConfig #1: fatal error config file: ", err)
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
		log.Error("GenerateUUID #0: ", err)
		return "", err
	}
	return id.String(), nil
}

// GenerateToken генерирует токен
func GenerateToken(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Error("GenerateToken #0: ", err)
		return "", err
	}

	token := hex.EncodeToString(randomBytes)
	return token, nil
}

// GenerateTokenBase64 генерирует токен в base64
func GenerateTokenBase64(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Error("GenerateToken #0: ", err)
		return "", err
	}

	// Кодирование случайной строки в base64
	token := base64.URLEncoding.EncodeToString(randomBytes)
	return token, nil
}

// GetHostIP получить все ip адреса на хосте
func GetHostIP() ([]net.IP, error) {
	host, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("GetHostIP #0: ", err)
		return nil, err
	}

	var ips []net.IP
	for _, addr := range host {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}

// RemovePrefixAuthBearer удаление префикса базовой авторизации для выделения токена
func RemovePrefixAuthBearer(t string) string {
	pref := "Bearer "
	if strings.HasPrefix(t, pref) {
		return strings.TrimPrefix(t, pref)
	}

	return t
}
