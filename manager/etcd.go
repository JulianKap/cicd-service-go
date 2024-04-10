package manager

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// setMasterETCD добавляет мастера в etcd
func (m *Master) setMasterETCD(cli *clientv3.Client) error {
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

	_, err = cli.Put(ctx, m.KEY, string(masterJSON), clientv3.WithLease(resp.ID))
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
func (w *Worker) setWorkerETCD(cli *clientv3.Client) error {
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

	_, err = cli.Put(ctx, w.KEY, string(workerJSON), clientv3.WithLease(resp.ID))
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

//// setConfig добавляет конфигурацию в etcd
//func setConfig(cli *clientv3.Client, workerID string) error {
//	// Контекст для операций
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Преобразование структуры Worker в JSON
//	workerJSON, err := json.Marshal(w)
//	if err != nil {
//		return err
//	}
//
//	// Добавление данных в etcd
//	_, err = cli.Put(ctx, fmt.Sprintf("/service-config/manager/workers/%s", workerID), string(workerJSON))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func isTTLValid(cli *clientv3.Client, key string) (bool, error) {
	ctx := context.Background()

	// Получение информации о ключе (включая TTL)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("isTTLValid error #0: ", err, ", key: ", key)
		return false, err
	}

	// Проверка наличия ключа и его TTL
	if len(resp.Kvs) == 0 {
		log.Info("isTTLValid info #1: key ", key, " not found")
		return false, nil // Ключ не найден
	}

	// Извлечение лиза из ответа
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

//func updateTTL(cli *clientv3.Client, key string, value string, ttl int64) error {
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
