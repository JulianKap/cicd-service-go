package vault

import (
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// SetToken добавление токена
func SetToken(cli *api.Client, t *Token) error {
	if _, err := cli.Logical().Write(t.Path, map[string]interface{}{"token": t.Token}); err != nil {
		log.Error("SetToken #0: ", err)
		return err
	}

	return nil
}

// GetToken  получить токен
func GetToken(cli *api.Client, t *Token) error {
	secret, err := cli.Logical().Read(t.Path)
	if err != nil {
		log.Error("GetToken #0: ", err)
		return err
	}

	if secret == nil {
		log.Error("GetToken #1: token not found ", t.Path)
	}

	token, ok := secret.Data["token"].(string)
	if !ok {
		log.Error("GetToken #2: ", err)
		return err
	}
	t.Token = token

	return nil
}

// DelToken удалить токен
func DelToken(cli *api.Client, t *Token) error {
	if _, err := cli.Logical().Delete(t.Path); err != nil {
		log.Error("DelToken #0: ", err)
		return err
	}

	return nil
}
