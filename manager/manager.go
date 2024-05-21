package manager

import (
	"cicd-service-go/constants"
	"cicd-service-go/utility"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Conf             Config  // Конфигурация кластера
	MemberInfo       Member  // Текущий экземпляр сервисы
	Keys             KeysDCS // Ключи в DCS для управления кластером
	CloseManagerChan chan bool
)

func InitManager() {
	Conf = Config{
		Cluster: ClusterConfig{
			Namespace:    viper.GetString("cluster.namespace_dcs"),
			TTL:          viper.GetInt("cluster.ttl"),
			LoopWait:     viper.GetUint("cluster.loop_wait"),
			RetryTimeout: viper.GetUint("cluster.retry_timeout"),
		},
	}

	Keys = KeysDCS{
		Config:  Conf.Cluster.Namespace + constants.CONFIG,
		Master:  Conf.Cluster.Namespace + constants.MASTER,
		Workers: Conf.Cluster.Namespace + constants.WORKERS,
		Members: Conf.Cluster.Namespace + constants.MEMBERS,
	}

	uniqueID, err := utility.GenerateUUID()
	if err != nil {
		log.Fatal("InitManager #0: ", err)
		return
	}
	MemberInfo.UUID = uniqueID
	Keys.Worker = Keys.Workers + "/" + uniqueID

	CloseManagerChan = make(chan bool)
}

func RunManager() {
	log.Info("RunManager #0: running UUID ", MemberInfo.UUID)

	clusterTicker := time.NewTicker(time.Duration(Conf.Cluster.RetryTimeout) * time.Second)
	for {
		select {
		case <-CloseManagerChan:
			log.Info("RunManager #0: close checking cluster state")
			return
		case <-clusterTicker.C:
			if err := tasksCluster(); err != nil {
				log.Error("RunManager #1: ", err)
			}
		}
	}
}

func tasksCluster() error {
	//// Инициализация кластера
	//if err := Conf.initializeCluster(); err != nil {
	//	log.Error("tasksCluster #0: ", err)
	//	return err
	//}

	// Проверяем конфигурацию сервиса в etcd
	if err := Conf.applyConfigurations(); err != nil {
		log.Error("tasksCluster #1: ", err)
		return err
	}

	// Проверяем наличие актуального мастера
	state, err := Conf.Cluster.checkMaster(&MemberInfo)
	if err != nil {
		log.Error("tasksCluster #2: ", err)
		return err
	}

	if state.Unknown {
		MemberInfo.ReadOnly = true
		log.Error("tasksCluster #3: ", err)
		return nil
	} else if !state.Exists || state.IAmMaster { // Становимся (или обновляем состояние) мастером
		if err := Conf.Cluster.setMaster(&MemberInfo); err != nil {
			log.Error("tasksCluster #4: ", err)
			return err
		}
		MemberInfo.Role = MasterRole

		if !state.IAmMaster {
			log.Info("tasksCluster #5: i became a MASTER!")
		}
	} else { // становимся слейвом
		MemberInfo.Role = WorkerRole
		log.Debug("tasksCluster #6: i became a SLAVE!")
	}

	if err := Conf.Cluster.setWorker(&MemberInfo); err != nil {
		log.Error("tasksCluster #7: ", err)
		return err
	}

	// Обновляем список членов кластера и удаляем старые (при просроченном TTL)
	if err := Conf.Cluster.updateMembers(&MemberInfo); err != nil {
		log.Error("tasksCluster #8: ", err)
		return err
	}

	MemberInfo.ReadOnly = false

	return nil
}
