package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Настройки подключения к etcd
type etcdConnectionOptions struct {
	endpoints   []string
	dialTimeout int
}

// InstanceETCD ..
var InstanceETCD *clientv3.Client

// InitInstanceETCD ..
func InitInstanceETCD() {
	InstanceETCD = initETCD(etcdConnectionOptions{
		endpoints:   viper.GetStringSlice("database.etcd.endpoints"),
		dialTimeout: viper.GetInt("database.etcd.dialTimeout"),
	})
}

func initETCD(options etcdConnectionOptions) *clientv3.Client {
	log.Debug("initETCD #0: start init etcd")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   options.endpoints,
		DialTimeout: time.Duration(options.dialTimeout) * time.Second,
	})
	if err != nil {
		log.Fatal("initETCD #1: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = cli.Get(ctx, "ping_key") // todo: переделать проверку пинга к etcd
	if err != nil {
		log.Fatal("initETCD #2: ", err)
	}

	log.Debug("initETCD #3: init etcd done")

	return cli
}
