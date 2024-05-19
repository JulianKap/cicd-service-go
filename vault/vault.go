package vault

import (
	"errors"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// SetToken добавление токена
func SetToken(cli *api.Client, t *Token) error {
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"token": t.Token,
		},
	}

	if _, err := cli.Logical().Write(t.Path, data); err != nil {
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

	if secret == nil || secret.Data == nil {
		log.Error("GetToken #1: token not found ", t.Path)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Error("GetToken #2: invalid data format at path ", t.Path)
		return errors.New("Invalid data format at path: " + t.Path)
	}

	token, ok := data["token"].(string)
	if !ok {
		log.Println("GetToken #3: not found at path ", t.Path)
		return errors.New("Not found at path: " + t.Path)
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
