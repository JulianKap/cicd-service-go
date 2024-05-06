package etcd

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// SetKey добавление ключа в etcd
func SetKey(cli *clientv3.Client, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cli.Put(ctx, key, value); err != nil {
		log.Error("SetKey #0: ", err)
		return err
	}

	return nil
}

// SetKeyTTL добавление ключа в etcd с ttl
func SetKeyTTL(cli *clientv3.Client, key string, value string, ttl int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Обновление TTL ключа
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		log.Error("SetKeyTTL #0: ", err)
		return err
	}

	_, err = cli.Put(ctx, key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Error("SetKeyTTL #1: ", err)
		return err
	}

	_, err = cli.KeepAliveOnce(ctx, resp.ID)
	if err != nil {
		log.Error("setKeyTTL #2: ", err)
		return err
	}

	return nil
}

// IsTTLValid проверка актуальности ключа по ttl
func IsTTLValid(cli *clientv3.Client, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение информации о ключе (включая TTL)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("IsTTLValid #0: ", err, ", key: ", key)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 {
		log.Info("IsTTLValid #1: key ", key, " not found")
		return false, nil
	}

	// Извлечение Lease из ответа
	leaseID := resp.Kvs[0].Lease

	// Получение TTL для ключа
	ttlResp, err := cli.TimeToLive(ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		log.Error("IsTTLValid #2: ", err, ", key: ", key)
		return false, err
	}

	if ttlResp.TTL == -1 {
		log.Info("IsTTLValid #3: key ", key, " has no TTL (does not expire)")
		return true, nil
	} else if ttlResp.TTL > 0 {
		log.Debug("IsTTLValid #4: TTL key ", key, " is still relevant")
		return true, nil
	}

	// TTL ключа истек
	log.Debug("IsTTLValid #5: key ", key, " TTL has expired")

	return false, nil
}

// GetKey получить значение ключа из etcd
func GetKey(cli *clientv3.Client, key string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Error("GetKey #0: ", err)
		return nil, err
	}

	return resp, nil
}

// DelKey удалить ключ из etcd
func DelKey(cli *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cli.Delete(ctx, key); err != nil {
		log.Error("DelKey #0: ", err)
		return err
	}

	return nil
}
