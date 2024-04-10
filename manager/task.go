package manager

import (
	"cicd-service-go/constants"
	"cicd-service-go/init/db"
	log "github.com/sirupsen/logrus"
)

// initializeCluster инициализируем кластер
func (c *Config) initializeCluster() error {
	return nil
}

// applyConfigurations проверяем конфигурацию в DCS
// Если отсутсвует, то добавляем.
func (c *Config) applyConfigurations() error {
	return nil
}

// checkMaster проверяем состояние мастера
func (m *Member) checkMaster(c *ClusterConfig) (MasterState, error) {
	state := MasterState{
		EXISTS:    false,
		INDEFINED: false,
	}

	key, err := isTTLValid(db.InstanceETCD, c.NAMESPACE+constants.MASTER)
	if err != nil {
		log.Error("setMaster error #0: ", err)
		state.INDEFINED = true
		return state, err
	}

	state.EXISTS = key

	return state, nil
}

// setMaster становимся мастером
func (m *Member) setMaster(c *ClusterConfig) error {
	master := Master{
		UUID:    m.UUID,
		TTL:     c.TTL,
		KEY:     c.NAMESPACE + constants.MASTER,
		RUNNING: true,
	}

	if err := master.setMasterETCD(db.InstanceETCD); err != nil {
		log.Error("setMaster error #0: ", err)
		return err
	}

	log.Debug("setMaster debug #1: I became a MASTER!")

	return nil
}

// setSlave становимся слейвом
func (m *Member) setSlave(c *ClusterConfig) error {
	worker := Worker{
		UUID:    m.UUID,
		TTL:     c.TTL,
		KEY:     c.NAMESPACE + constants.WORKERS + "/" + m.UUID,
		RUNNING: true,
	}

	if err := worker.setWorkerETCD(db.InstanceETCD); err != nil {
		log.Error("setSlave error #0: ", err)
		return err
	}

	log.Debug("setSlave debug #1: I became a SLAVE!")

	return nil
}
