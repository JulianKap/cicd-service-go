package manager

// ClusterConfig основные настройки кластера
type ClusterConfig struct {
	NAMESPACE     string
	TTL           int
	LOOP_WAIT     uint
	RETRY_TIMEOUT uint
}

// Config общий конфиг сервиса
type Config struct {
	CLUSTER ClusterConfig
}

// Member текущее состояние данного экземпляра сервиса
type Member struct {
	UUID       string `json:"uuid"`
	Master     bool   `json:"master"`
	STANDALONE bool   `json:"standalone"` // когда только один узел
}

// MasterState текущее состояние мастера
type MasterState struct {
	EXISTS    bool `json:"exists"`
	INDEFINED bool `json:"indefined"`
}

// Master информация о мастере в etcd
type Master struct {
	UUID    string `json:"uuid"`
	TTL     int    `json:"ttl"`
	KEY     string `json:"key"`
	RUNNING bool   `json:"running"`
	//Standalone bool   `json:"standalone"`
}

// Worker информация о рабочем узле в etcd
type Worker struct {
	UUID    string `json:"uuid"`
	TTL     int    `json:"ttl"`
	KEY     string `json:"key"`
	RUNNING bool   `json:"running"`
}
