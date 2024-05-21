package secrets

import (
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type vaultConnectionOptions struct {
	addr    string
	token   string
	timeout int
}

var InstanceVault *api.Client

func InitInstanceVault() {
	InstanceVault = initVault(vaultConnectionOptions{
		addr:    viper.GetString("database.vault.addr"),
		token:   viper.GetString("database.vault.token"),
		timeout: viper.GetInt("database.vault.timeout_s"),
	})
}

func initVault(options vaultConnectionOptions) *api.Client {
	log.Debug("initVault #0: start init secrets")

	config := &api.Config{
		Address: options.addr,
		Timeout: time.Duration(options.timeout) * time.Second,
	}
	cli, err := api.NewClient(config)
	if err != nil {
		log.Fatal("initVault #1: ", err)
	}

	if options.token == "" {
		log.Fatal("initVault #2: token is null")
	}
	cli.SetToken(options.token)

	// Проверка подключения к Vault
	if _, err := cli.Sys().Health(); err != nil {
		log.Fatal("initVault #3: ", err)
	}

	//todo: сделать проверку, что токен действительный через тестовое создание и чтение ключей через отдельную функцию

	log.Debug("initVault #4: init secrets done")

	return cli
}
