package manager

import (
	"cicd-service-go/db/etcd"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// setConfigETCD добавление конфигурации в etcd
func (c *Config) setConfigETCD(cli *clientv3.Client) error {
	valueJSON, err := json.Marshal(c)
	if err != nil {
		log.Error("setConfigETCD #0: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Config, string(valueJSON)); err != nil {
		log.Error("setConfigETCD #1: ", err)
		return err
	}

	return nil
}

// setMasterETCD добавляет мастера в etcd
func (m *Master) setMasterETCD(cli *clientv3.Client) error {
	valueJSON, err := json.Marshal(m)
	if err != nil {
		log.Error("setMasterETCD #0: ", err)
		return err
	}

	if err = etcd.SetKeyTTL(cli, Keys.Master, string(valueJSON), m.TTL); err != nil {
		log.Error("setMasterETCD #1: ", err)
		return err
	}

	return nil
}

// setWorkerETCD добавляет рабочего в etcd
func (w *Worker) setWorkerETCD(cli *clientv3.Client) error {
	valueJSON, err := json.Marshal(w)
	if err != nil {
		log.Error("setWorkerETCD #0: ", err)
		return err
	}

	if err = etcd.SetKeyTTL(cli, Keys.Worker, string(valueJSON), w.TTL); err != nil {
		log.Error("setWorkerETCD #1: ", err)
		return err
	}

	return nil
}

// setMembersETCD добавляет в список членов кластера
func (m *Members) setMembersETCD(cli *clientv3.Client, key string) error {
	valueJSON, err := json.Marshal(m)
	if err != nil {
		log.Error("setMembersETCD #0: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Members, string(valueJSON)); err != nil {
		log.Error("setMembersETCD #1: ", err)
		return err
	}

	return nil
}

// getMasterETCD получение мастера из etcd
func getMasterETCD(cli *clientv3.Client) (*Master, error) {
	resp, err := etcd.GetKey(cli, Keys.Master)
	if err != nil {
		log.Error("getMasterETCD #0: ", err)
		return nil, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getMasterETCD #1: key ", Keys.Master, " not found")
		return nil, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getMasterETCD #2: key ", Keys.Master, " get more than one key")
	}

	master := Master{}
	if err := json.Unmarshal(resp.Kvs[0].Value, &master); err != nil {
		log.Error("getMasterETCD #3: ", err)
		return nil, err
	}

	return &master, nil
}

func getMembersETCD(cli *clientv3.Client, key string) (*Members, error) {
	resp, err := etcd.GetKey(cli, Keys.Members)
	if err != nil {
		log.Error("getMembersETCD #0: ", err)
		return nil, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getMembersETCD #1: key ", key, " not found")
		return nil, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getMembersETCD #2: key ", key, " get more than one key")
	}

	members := Members{}
	if err := json.Unmarshal(resp.Kvs[0].Value, &members); err != nil {
		log.Error("getMembersETCD #3: ", err)
		return nil, err
	}

	return &members, nil
}

// delMemberETCD удалить члена кластера в etcd
func (m *Member) delMemberETCD(cli *clientv3.Client, key string) error {
	if err := etcd.DelKey(cli, key); err != nil {
		log.Error("delMemberETCD #0: ", err)
		return err
	}

	return nil
}
