package manager

// ClusterConfig основные настройки кластера
type ClusterConfig struct {
	Namespace    string
	TTL          int
	LoopWait     uint
	RetryTimeout uint
}

// Config общий конфиг сервиса
type Config struct {
	Cluster ClusterConfig
}

// Member текущее состояние данного экземпляра сервиса
type Member struct {
	UUID       string `json:"uuid"`
	Master     bool   `json:"master"`
	Standalone bool   `json:"standalone"` // когда только один узел
}

// Members список всех членов кластера
type Members struct {
	Members []Member `json:"members"`
}

// MasterState текущее состояние мастера
type MasterState struct {
	Exists    bool `json:"exists"`
	IAmMaster bool `json:"i_am_master"`
	Unknown   bool `json:"unknown"` // неопределенный статус, вероятно ошибка с etcd. Следует полождать
}

// Master информация о мастере в etcd
type Master struct {
	UUID    string `json:"uuid"`
	TTL     int    `json:"ttl"`
	Running bool   `json:"running"`
	//Standalone bool   `json:"standalone"`
}

// Worker информация о рабочем узле в etcd
type Worker struct {
	UUID    string `json:"uuid"`
	TTL     int    `json:"ttl"`
	Running bool   `json:"running"`
}
