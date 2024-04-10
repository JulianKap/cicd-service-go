package constants

const (
	PROJECTS    = "/projects"
	MASTER      = "/master"
	WORKERS     = "/workers"
	STATUS      = "/status"
	CONFIG      = "/config"
	MASTER_UUID = ""
	MASTER_
)

//
//// Базовые ключи в etcd
//var KeysBase = struct {
//	PROJECTS string // Каталог с проектами и задачами
//	WORKERS  string // члены кластер
//	MASTER   string // мастер
//	STATUS   string // статус кластера
//	CONFIG   string // конфигурация кластера
//}{
//	PROJECTS: "/projects",
//	WORKERS:  "/workers",
//	MASTER:   "/master",
//	STATUS:   "/status",
//	CONFIG:   "/config",
//}
//
//// Ключи проектов
//var KeysProject = struct {
//	APIKey string
//	NAME   string
//	JOBS   string
//}{
//	APIKey: "/leader",
//	NAME:   "/status",
//	JOBS:   "",
//}
//
//// Ключи задач
//var KeysJob = struct {
//	NAME   string
//	URL    string
//	Branch string
//}{
//	NAME:   "/leader",
//	URL:    "/config",
//	Branch: "",
//}
