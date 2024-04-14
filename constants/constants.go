package constants

const (
	// Структура etcd
	PROJECTS          = "/projects"
	PROJECTS_ALL      = "/projects/all"
	PROJECT_LATEST_ID = "/projects/latest_id"
	JOBS              = "/jobs"
	JOB_LATEST_ID     = "/jobs/latest_id"
	MASTER            = "/master"
	WORKERS           = "/workers"
	MEMBERS           = "/members"
	//STATUS   = "/status"
	CONFIG = "/config"
	// Настройки подключения
	CONTEXT_TIMEOUT_ETCD = 5 // todo: добавить в контексте
)
