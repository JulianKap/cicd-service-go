package manager

import (
	"cicd-service-go/db/etcd"
	"cicd-service-go/init/db"

	log "github.com/sirupsen/logrus"
)

//// initializeCluster инициализация кластера
//func (c *Config) initializeCluster() error {
//	return nil
//}

// applyConfigurations проверка конфигурации в DCS
// Если конфигурация отсутсвует, то добавляем ее
func (c *Config) applyConfigurations() error {
	if err := c.setConfigETCD(db.InstanceETCD); err != nil {
		log.Error("applyConfigurations #0: ", err)
		return err
	}
	log.Debug("applyConfigurations #1: update configuration")

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
		log.Error("checkMaster #0: ", err)
		state.Unknown = true
		return state, err
	}

	// Поскольку ключ валидный, проверим совпадение ключа текущего экземпляра сервиса с мастером
	if validate {
		state.Exists = true

		res, err := getMasterETCD(db.InstanceETCD)
		if err != nil {
			log.Error("checkMaster #1: ", err)
			return state, err
		}

		if res.UUID == m.UUID {
			log.Debug("checkMaster #2: i am a MASTER")
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
		log.Error("setMaster #0: ", err)
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
		log.Error("setWorker #0: ", err)
		return err
	}

	return nil
}

// updateMembers добавление новой реплики в список кластера.
// Также удаление старых реплики, у которых просрочен ttl
func (c *ClusterConfig) updateMembers(m *Member) error {
	res, err := GetMembers()
	if err != nil {
		log.Error("updateMembers #0: ", err)
		return err
	}

	newMembers := Members{}
	newMembers.Members = append(newMembers.Members, *m)

	delMembers := Members{}

	if res != nil {
		if len(res.Members) > 0 {
			if m.Role == MasterRole {
				for _, mbr := range res.Members {
					if mbr.UUID != m.UUID {
						keyWorker := Keys.Workers + "/" + mbr.UUID

						validate, err := etcd.IsTTLValid(db.InstanceETCD, keyWorker)
						if err != nil {
							log.Error("updateMembers #1: ", err)
							continue
						}

						if validate {
							newMembers.Members = append(newMembers.Members, mbr) // Актуальные члены члены кластера
						} else {
							delMembers.Members = append(delMembers.Members, mbr) // Не актуальные
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
	}

	if err := newMembers.setMembersETCD(db.InstanceETCD); err != nil {
		log.Error("updateMembers #3: ", err)
		return err
	}

	//todo:сделать реестр старых членов кластера с лимитом на кол-во элементов
	//if err := delMembers.delMembersETCD(db.InstanceETCD); err != nil {
	//	log.Error("updateMembers #4: ", err)
	//	return err
	//}

	return nil
}

// GetMembers получение списка членов кластера
func GetMembers() (*Members, error) {
	res, err := getMembersETCD(db.InstanceETCD, Keys.Members)
	if err != nil {
		log.Error("GetMembers #0: ", err)
		return nil, err
	}

	return res, nil
}
