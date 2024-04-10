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
	key := c.Cluster.Namespace + constants.CONFIG

	if err := c.setConfigETCD(db.InstanceETCD, key); err != nil {
		log.Error("applyConfigurations error #0: ", err)
		return err
	}

	log.Debug("applyConfigurations debug #1: update configuration")

	return nil
}

// checkMaster проверяем состояние мастера
func (m *Member) checkMaster(c *ClusterConfig) (MasterState, error) {
	state := MasterState{
		Exists:    false,
		IAmMaster: false,
		Unknown:   false,
	}

	validate, err := isTTLValid(db.InstanceETCD, c.Namespace+constants.MASTER)
	if err != nil {
		log.Error("checkMaster error #0: ", err)
		state.Unknown = true
		return state, err
	}
	state.Exists = validate

	if validate { //Поскольку ключ валидный, проверим совпадение ключа
		res, err := getMasterETCD(db.InstanceETCD, c.Namespace+constants.MASTER)
		if err != nil {
			log.Error("setMaster error #1: ", err)
			state.Unknown = true
			return state, err
		}

		if res.UUID == m.UUID {
			state.IAmMaster = true
		}
	}

	return state, nil
}

// setMaster становимся мастером
func (m *Member) setMaster(c *ClusterConfig) error {
	master := Master{
		UUID:    m.UUID,
		TTL:     c.TTL,
		Running: true,
	}

	key := c.Namespace + constants.MASTER

	if err := master.setMasterETCD(db.InstanceETCD, key); err != nil {
		log.Error("setMaster error #0: ", err)
		return err
	}
	m.Master = true

	if err := m.updateMembers(c); err != nil {
		log.Error("setMaster error #1: ", err)
		return err
	}

	log.Debug("setMaster debug #2: I became a MASTER!")

	return nil
}

// setSlave становимся слейвом
func (m *Member) setSlave(c *ClusterConfig) error {
	worker := Worker{
		UUID:    m.UUID,
		TTL:     c.TTL,
		Running: true,
	}

	key := c.Namespace + constants.WORKERS + "/" + m.UUID

	if err := worker.setWorkerETCD(db.InstanceETCD, key); err != nil {
		log.Error("setSlave error #0: ", err)
		return err
	}
	m.Master = false

	if err := m.updateMembers(c); err != nil {
		log.Error("setSlave error #1: ", err)
		return err
	}

	log.Debug("setSlave debug #2: I became a SLAVE!")

	return nil
}

// updateMembers добавление инстанса в список кластера
func (m *Member) updateMembers(c *ClusterConfig) error {
	key := c.Namespace + constants.MEMBERS

	res, err := getMembersETCD(db.InstanceETCD, key)
	if err != nil {
		log.Error("updateMembers error #0: ", err)
		return err
	}

	for _, memb := range res.Members {
		if memb.UUID == m.UUID {
			return nil // todo: сделать обновление полей для члена
		}
	}

	res.Members = append(res.Members, *m)

	if err := res.setMembersETCD(db.InstanceETCD, key); err != nil {
		log.Error("updateMembers error #1: ", err)
		return err
	}

	return nil
}
