package manager

import (
	"cicd-service-go/constants"
	"cicd-service-go/utility"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	config      Config  // конфигурация
	MemberInfo  Member  // состояние текущего экзмпляра
	Keys        KeysDCS // ключи в DCS для управления кластером
	managerChan chan bool
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

	Keys = KeysDCS{
		Config:  config.Cluster.Namespace + constants.CONFIG,
		Master:  config.Cluster.Namespace + constants.MASTER,
		Workers: config.Cluster.Namespace + constants.WORKERS,
		Members: config.Cluster.Namespace + constants.MEMBERS,
	}

	uniqueID, err := utility.GenerateUUID()
	if err != nil {
		log.Fatal("InitManager error #0: ", err)
		return
	}
	MemberInfo.UUID = uniqueID
	Keys.Worker = Keys.Workers + "/" + uniqueID

	managerChan = make(chan bool)
}

//func GetMemberChan() chan bool {
//	return managerChan
//}

func RunManager() {
	log.Info("RunManager info #0: running UUID ", MemberInfo.UUID)

	clusterTicker := time.NewTicker(time.Duration(config.Cluster.RetryTimeout) * time.Second)
	for {
		select {
		case <-managerChan:
			log.Info("RunManager info #0: close checking cluster state")
			return
		case <-clusterTicker.C:
			if err := tasksCluster(); err != nil {
				log.Error("RunManager error #1: ", err)
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
	state, err := config.Cluster.checkMaster(&MemberInfo)
	if err != nil {
		log.Error("tasksCluster error #2: ", err)
		return err
	}

	if state.Unknown {
		MemberInfo.ReadOnly = true
		log.Error("tasksCluster error #3: ", err)
		return nil
	} else if !state.Exists || state.IAmMaster { // Становимся (или обновляем состояние) мастером
		if err := config.Cluster.setMaster(&MemberInfo); err != nil {
			log.Error("tasksCluster error #4: ", err)
			return err
		}
		MemberInfo.Master = true

		if !state.IAmMaster {
			log.Info("tasksCluster info #5: i became a MASTER!")
		}
	} else { // становимся слейвом
		MemberInfo.Master = false
		log.Debug("tasksCluster debug #6: i became a SLAVE!")
	}

	if err := config.Cluster.setWorker(&MemberInfo); err != nil {
		log.Error("tasksCluster error #7: ", err)
		return err
	}

	// Обновляем список членов кластера и удаляем старые (при просроченном TTL)
	if err := config.Cluster.updateMembers(&MemberInfo); err != nil {
		log.Error("tasksCluster error #8: ", err)
		return err
	}

	// Проверять таски для каждого воркера перед удалением. Либо удалять, а в шедулере будет своя логика распределения незаконченных задач

	MemberInfo.ReadOnly = false

	return nil
}
