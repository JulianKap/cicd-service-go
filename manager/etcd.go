package manager

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// setConfigETCD добавляет конфигурацию в etcd
func (c *Config) setConfigETCD(cli *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	workerJSON, err := json.Marshal(c)
	if err != nil {
		log.Error("setConfigETCD error #0: ", err)
		return err
	}

	if _, err = cli.Put(ctx, key, string(workerJSON)); err != nil {
		log.Error("setConfigETCD error #1: ", err)
		return err
	}

	return nil
}

// setMasterETCD добавляет мастера в etcd
func (m *Master) setMasterETCD(cli *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	masterJSON, err := json.Marshal(m)
	if err != nil {
		log.Error("setMasterETCD error #0: ", err)
		return err
	}

	// Обновление TTL ключа
	resp, err := cli.Grant(ctx, int64(m.TTL))
	if err != nil {
		log.Error("setMasterETCD error #1: ", err)
		return err
	}

	_, err = cli.Put(ctx, key, string(masterJSON), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Error("setMasterETCD error #2: ", err)
		return err
	}

	_, err = cli.KeepAliveOnce(ctx, resp.ID)
	if err != nil {
		log.Error("setMasterETCD error #3: ", err)
		return err
	}

	return nil
}

// setWorkerETCD добавляет рабочего в etcd
func (w *Worker) setWorkerETCD(cli *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	workerJSON, err := json.Marshal(w)
	if err != nil {
		log.Error("setWorkerETCD error #0: ", err)
		return err
	}

	// Обновление TTL ключа
	resp, err := cli.Grant(ctx, int64(w.TTL))
	if err != nil {
		log.Error("setWorkerETCD error #1: ", err)
		return err
	}

	_, err = cli.Put(ctx, key, string(workerJSON), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Error("setWorkerETCD error #2: ", err)
		return err
	}

	_, err = cli.KeepAliveOnce(ctx, resp.ID)
	if err != nil {
		log.Error("setWorkerETCD error #3: ", err)
		return err
	}

	return nil
}

func isTTLValid(cli *clientv3.Client, key string) (bool, error) {
	ctx := context.Background()

	// Получение информации о ключе (включая TTL)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("isTTLValid error #0: ", err, ", key: ", key)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 {
		log.Info("isTTLValid info #1: key ", key, " not found")
		return false, nil // Ключ не найден
	}

	// Извлечение Lease из ответа
	leaseID := resp.Kvs[0].Lease

	// Получение TTL для ключа
	ttlResp, err := cli.TimeToLive(ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		log.Error("isTTLValid error #2: ", err, ", key: ", key)
		return false, err
	}

	// Проверка TTL
	if ttlResp.TTL == -1 {
		log.Info("isTTLValid info #3: key ", key, " has no TTL (does not expire)")
		return true, nil
	} else if ttlResp.TTL > 0 {
		log.Debug("isTTLValid debug #4: TTL key ", key, " is still relevant")
		return true, nil
	}

	// TTL ключа истек
	log.Debug("isTTLValid debug #5: key ", key, " TTL has expired")

	return false, nil
}

// getMasterETCD получение мастера из etcd
func getMasterETCD(cli *clientv3.Client, key string) (Master, error) {
	master := Master{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение значения ключа из etcd
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("getMasterETCD error #0: ", err)
		return master, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getMasterETCD info #1: key ", key, " not found")
		return master, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Error("getMasterETCD error #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &master); err != nil {
		log.Error("getMasterETCD error #3: ", err)
		return master, err
	}

	return master, nil
}

//func setKeyWithTTL(cli *clientv3.Client, key string, value string, ttl int64) error {
//	// Контекст для отмены операций
//	ctx := context.Background()
//
//	// Обновление TTL ключа
//	resp, err := cli.Grant(ctx, ttl)
//	if err != nil {
//		return err
//	}
//
//	// Установка нового TTL для ключа
//	_, err = cli.Put(ctx, key, value, clientv3.WithLease(resp.ID))
//	if err != nil {
//		return err
//	}
//
//	_, err = cli.KeepAliveOnce(ctx, resp.ID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (m *Members) setMembersETCD(cli *clientv3.Client, key string) error {
	ctx := context.Background()

	workerJSON, err := json.Marshal(m)
	if err != nil {
		log.Error("setMembersETCD error #0: ", err)
		return err
	}

	_, err = cli.Put(ctx, key, string(workerJSON))
	if err != nil {
		log.Error("setMembersETCD error #1: ", err)
		return err
	}

	return nil
}

func getMembersETCD(cli *clientv3.Client, key string) (Members, error) {
	members := Members{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение значения ключа из etcd
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("getMembersETCD error #0: ", err)
		return members, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getMembersETCD info #1: key ", key, " not found")
		return members, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Error("getMembersETCD error #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &members); err != nil {
		log.Error("getMembersETCD error #3: ", err)
		return members, err
	}

	return members, nil
}
