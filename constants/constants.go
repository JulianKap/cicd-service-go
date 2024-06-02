package constants

const (
	// Ключи DCS с абсолютным путем
	PROJECTS                 = "/projects"
	PROJECTS_ALL             = "/projects/all"
	PROJECT_LATEST_ID        = "/projects/latest_id"
	PROJECTS_TASKS           = "/projects/tasks"
	PROJECTS_TASKS_HISTORY   = "/projects/tasks_history"
	PROJECTS_TASKS_LATEST_ID = "/projects/tasks/latest_id"
	MASTER                   = "/master"
	WORKERS                  = "/workers"
	MEMBERS                  = "/members"
	STATUS                   = "/status"
	CONFIG                   = "/config"

	// Ключи DCS с относительным путем
	JOBS          = "/jobs"
	JOB_LATEST_ID = "/jobs/latest_id"
	TASKS         = "/tasks"

	// Настройки подключения
	CONTEXT_TIMEOUT_ETCD = 5 // todo: добавить в контексте
	// todo: добавить длину токена генериации, чтоб было 16 символов
)
