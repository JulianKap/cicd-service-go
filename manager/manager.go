package manager

import (
	"cicd-service-go/utility"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	config     Config
	member     Member
	memberChan chan bool
)

func InitManager() {
	config = Config{
		CLUSTER: ClusterConfig{
			NAMESPACE:     viper.GetString("cluster.namespace_dcs"),
			TTL:           viper.GetInt("cluster.ttl"),
			LOOP_WAIT:     viper.GetUint("cluster.loop_wait"),
			RETRY_TIMEOUT: viper.GetUint("cluster.retry_timeout"),
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

	member.UUID = uniqueID
	log.Debug("Generated new UUID: ", member.UUID)

	clusterTicker := time.NewTicker(time.Duration(config.CLUSTER.RETRY_TIMEOUT) * time.Second)

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
	state, err := member.checkMaster(&config.CLUSTER)
	if err != nil {
		log.Error("tasksCluster error #2: ", err)
		return err
	}

	if state.EXISTS { // Становимся воркером

		// Проверить UUID, чтоб совпадал
		// Проверить, что есть актуальные воркеры по ttl. Если нет, то переключаемся в режим стенделоне

		if err := member.setSlave(&config.CLUSTER); err != nil {
			log.Error("tasksCluster error #3: ", err)
			return err
		}
	} else { // Становимся мастером
		if err := member.setMaster(&config.CLUSTER); err != nil {
			log.Error("tasksCluster error #4: ", err)
			return err
		}
	}

	return nil
}
