package sources

import (
	"cicd-service-go/utility"
	"cicd-service-go/vault"
	"fmt"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// createProjectETCD создание нового проекта
func (p *Project) createTokenProjectVault(cli *api.Client, t *vault.Token) error {
	token, err := utility.GenerateToken(16)
	if err != nil {
		log.Error("createTokenProjectVault #0: ", err)
		return err
	}
	t.Token = token
	t.Path = fmt.Sprintf("project/%s/token", p.ProjectName)

	if err = vault.SetToken(cli, t); err != nil {
		log.Error("createTokenProjectVault #0: ", err)
		return err
	}

	return nil
}

func (p *Project) checkTokenVault(cli *api.Client, t *vault.Token) (bool, error) {
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
