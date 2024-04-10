package manager

import (
	"cicd-service-go/utility"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	config     Config
	MemberInfo Member
	memberChan chan bool
)

func InitManager() {
	config = Config{
		Cluster: ClusterConfig{
			Namespace:    viper.GetString("cluster.namespace_dcs"),
			TTL:          viper.GetInt("cluster.ttl"),
			LoopWait:     viper.GetUint("cluster.loop_wait"),
			RetryTimeout: viper.GetUint("cluster.retry_timeout"),
		},
	}

	memberChan = make(chan bool)
}

func GetMemberChan() chan bool {
	return memberChan
}

func RunManager() {
	// Генерируем уникальный UUID запущенного экземпляра сервиса.
	uniqueID, err := utility.GenerateUUID()
	if err != nil {
		log.Fatal("RunManager error #0: ", err)
		return
	}

	MemberInfo.UUID = uniqueID
	log.Debug("Generated new UUID: ", MemberInfo.UUID)

	clusterTicker := time.NewTicker(time.Duration(config.Cluster.RetryTimeout) * time.Second)

	for {
		select {
		case <-memberChan:
			log.Infoln("Close checking cluster state")
			return
		case <-clusterTicker.C:
			if err := tasksCluster(); err != nil {
				log.Debugln("RunManager error #1: ", err)
			}
		}
	}
}

func tasksCluster() error {
	// Инициализация кластера
	if err := config.initializeCluster(); err != nil {
		log.Error("tasksCluster error #0: ", err)
		return err
	}

	// Проверяем конфигурацию сервиса в etcd
	if err := config.applyConfigurations(); err != nil {
		log.Error("tasksCluster error #1: ", err)
		return err
	}

	// Проверяем наличие актуального мастера
	master, err := MemberInfo.checkMaster(&config.Cluster)
	if err != nil {
		log.Error("tasksCluster error #2: ", err)
		return err
	}

	if master.IAmMaster { //Мы являемся актуальным мастером

	} else if !master.Exists && !master.Unknown { // Становимся мастером
		if err := MemberInfo.setMaster(&config.Cluster); err != nil {
			log.Error("tasksCluster error #4: ", err)
			return err
		}
	} else if master.Exists && !master.Unknown { // Мастер существует в кластере, поэтому становимся воркером

		if err := MemberInfo.setSlave(&config.Cluster); err != nil {
			log.Error("tasksCluster error #3: ", err)
			return err
		}
	} else if master.Unknown {
		log.Error("tasksCluster error #4: ", err)
		return nil
	}

	// Проверить, что есть актуальные воркеры по ttl. Если нет, то переключаемся в режим стенделоне

	return nil
}
