package manager

import (
	"cicd-service-go/db/etcd"
	"cicd-service-go/init/db"
	log "github.com/sirupsen/logrus"
)

// initializeCluster инициализация кластера
func (c *Config) initializeCluster() error {
	return nil
}

// applyConfigurations проверка конфигурации в DCS
// Если конфигурация отсутсвует, то добавляем ее
func (c *Config) applyConfigurations() error {
	if err := c.setConfigETCD(db.InstanceETCD); err != nil {
		log.Error("applyConfigurations error #0: ", err)
		return err
	}
	log.Debug("applyConfigurations debug #1: update configuration")

	return nil
}

// checkMaster проверяем состояние мастера
func (c *ClusterConfig) checkMaster(m *Member) (MasterState, error) {
	state := MasterState{
		Exists:    false,
		IAmMaster: false,
		Unknown:   false,
	}

	validate, err := etcd.IsTTLValid(db.InstanceETCD, Keys.Master)
	if err != nil {
		log.Error("checkMaster error #0: ", err)
		state.Unknown = true
		return state, err
	}

	// Поскольку ключ валидный, проверим совпадение ключа текущего экземпляра сервиса с мастером
	if validate {
		state.Exists = true

		res, err := getMasterETCD(db.InstanceETCD)
		if err != nil {
			log.Error("checkMaster error #1: ", err)
			return state, err
		}

		if res.UUID == m.UUID {
			log.Debug("checkMaster debug #2: i am a MASTER")
			state.IAmMaster = true
		}
	}

	return state, nil
}

// setMaster становимся мастером
func (c *ClusterConfig) setMaster(m *Member) error {
	master := Master{
		UUID: m.UUID,
		TTL:  c.TTL,
	}

	if err := master.setMasterETCD(db.InstanceETCD); err != nil {
		log.Error("setMaster error #0: ", err)
		return err
	}

	return nil
}

// setWorker становимся воркером
// Стоит отметить, что мастер всегда является воркером, но воркер не всегда мастером
func (c *ClusterConfig) setWorker(m *Member) error {
	worker := Worker{
		UUID:    m.UUID,
		TTL:     c.TTL,
		Running: true,
	}

	if err := worker.setWorkerETCD(db.InstanceETCD); err != nil {
		log.Error("setWorker error #0: ", err)
		return err
	}

	return nil
}

// updateMembers добавление инстанса в список кластера.
// Также удаление старых инстансов, у которых просрочен ttl
func (c *ClusterConfig) updateMembers(m *Member) error {
	res, err := getMembersETCD(db.InstanceETCD, Keys.Members)
	if err != nil {
		log.Error("updateMembers error #0: ", err)
		return err
	}

	newMembers := Members{}
	newMembers.Members = append(newMembers.Members, *m)

	if res != nil {
		if m.Master {
			for _, mbr := range res.Members {
				if mbr.UUID != m.UUID {
					keyWorker := Keys.Workers + "/" + mbr.UUID

					validate, err := etcd.IsTTLValid(db.InstanceETCD, keyWorker)
					if err != nil {
						log.Error("updateMembers error #1: ", err)
						continue
					}

					if validate {
						newMembers.Members = append(newMembers.Members, mbr)
						// todo: как вариант, нужно будет удалять вручную ключи воркеров
					}
				}
			}
		} else {
			for _, mbr := range res.Members {
				if mbr.UUID != m.UUID {
					newMembers.Members = append(newMembers.Members, mbr)
				}
			}
		}
	}

	if err := newMembers.setMembersETCD(db.InstanceETCD, Keys.Members); err != nil {
		log.Error("updateMembers error #3: ", err)
		return err
	}

	return nil
}
