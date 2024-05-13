package manager

// KeysDCS ключи в DCS
type KeysDCS struct {
	// Config ключ конфига
	Config string `json:"config"`
	// Master ключ мастера
	Master string `json:"master"`
	// Worker ключ воркеров
	Workers string `json:"workers"`
	// Members ключ членов кластера
	Members string `json:"members"`
	// Worker ключ текущего воркера
	Worker string `json:"worker"`
}

// ClusterConfig конфиг кластераx
type ClusterConfig struct {
	Namespace    string `json:"namespace"`
	TTL          int    `json:"ttl"`
	LoopWait     uint   `json:"loop_wait"`
	RetryTimeout uint   `json:"retry_timeout"`
}

// Config общий конфиг сервиса
type Config struct {
	// Cluster настройки кластера
	Cluster ClusterConfig `json:"cluster"`
}

// RoleStatus роль, в которой сервис работает
type RoleStatus int

const (
	WorkerRole     RoleStatus = iota // Режим воркера
	MasterRole                       // Режим мастера
	StandaloneRole                   // Режим все в одном для одиночного разворачивания
)

// Member настройки данного экземпляра сервиса
type Member struct {
	// UUID уникальный id запущенного экземпляра сервиса
	UUID string `json:"uuid"`
	// Роль
	Role RoleStatus `json:"role"`
	// ReadOnly режим работы только для чтения (когда состояние кластера неизвестно)
	ReadOnly bool `json:"read_only"`
}

// Members список всех членов кластера
type Members struct {
	Members []Member `json:"members"`
}

// MasterState текущее состояние мастера в кластера
type MasterState struct {
	// Exists мастер существует
	Exists bool `json:"exists"`
	// IAmMaster данный экзмепляр является мастером
	IAmMaster bool `json:"i_am_master"`
	// Unknown неизвестный статус, вероятно ошибка с DCS. Следует через время повторить и перевестить в readonly режим
	Unknown bool `json:"unknown"`
}

// Master информация о мастере в etcd
type Master struct {
	UUID string `json:"uuid"`
	TTL  int    `json:"ttl"`
}

// Worker информация о рабочем узле в etcd
type Worker struct {
	UUID    string `json:"uuid"`
	TTL     int    `json:"ttl"`
	Running bool   `json:"running"`
	// Список тасок в работе. Указать параметр в конфиге с кол-вом хранимой истории с датой и статусом
}
