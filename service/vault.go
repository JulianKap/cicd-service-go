package service

import (
	"cicd-service-go/sources"
	"cicd-service-go/vault"
	"fmt"
	"github.com/hashicorp/vault/api"

	log "github.com/sirupsen/logrus"
)

// СheckTokenVault проверка токена авторизации
func СheckTokenVault(cli *api.Client, p *sources.Project, t *vault.Token) (bool, error) {
	var tVal vault.Token
	tVal.Path = fmt.Sprintf("project/%s/token", p.ProjectName)

	if err := vault.GetToken(cli, &tVal); err != nil {
		log.Error("checkTokenVault #0: ", err)
		return false, err
	}

	if t.Token != tVal.Token {
		return false, nil
	}

	return true, nil
}
