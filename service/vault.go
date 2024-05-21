package service

import (
	"cicd-service-go/sources"
	"cicd-service-go/vault"
	"github.com/hashicorp/vault/api"

	log "github.com/sirupsen/logrus"
)

// checkTokenVault проверка токена авторизации
func checkTokenVault(cli *api.Client, p *sources.Project, t *vault.Token) (bool, error) {
	var tVal vault.Token
	tVal.Path = sources.GetProjectPathToken(p)

	if err := vault.GetToken(cli, &tVal); err != nil {
		log.Error("checkTokenVault #0: ", err)
		return false, err
	}

	if t.Token != tVal.Token {
		return false, nil
	}

	return true, nil
}
