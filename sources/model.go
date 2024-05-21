package sources

// Credentials креды для авторизации
type Credentials struct {
	// Username пользователь
	Username string `json:"username,omitempty"`
	// Password пароль пользователя
	Password string `json:"-"`
	// Token токен авторизации
	Token string `json:"-"`
}

// Project представляет информацию о проекте
type Project struct {
	// ID уникальный id проекта
	ID int `json:"id,omitempty"`
	// ProjectName название проекта
	ProjectName string `json:"project_name"`
	// Token токен авторизации для проекта
	Token string `json:"token,omitempty"`
}

// Projects представляет список всех проектов
type Projects struct {
	Projects []*Project `json:"projects"`
}

// Job представляет информацию о задаче
type Job struct {
	// ID уникальный id задачи
	ID int `json:"id"`
	// ID id проекта, к которому относится задача
	IdProject int `json:"id_project"`
	// Название задачи
	JobName string `json:"job_name"`
	// URL url репозитория
	URL string `json:"url"`
	// Branch ветка репозитория
	Branch string `json:"branch"`
	// Enable активность задачи
	Enable bool `json:"enable"`
	// Creds креды для доступа к репозиторию todo: отдельную структуру для кредов
	Creds Credentials `json:"creds,omitempty"`
}

// Jobs представляет список всех задач проекта
type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

//// ProjectFull представляет расширенную информацию о проекте
//type ProjectFull struct {
//	// ID уникальный id проекта
//	ID int `json:"id,omitempty"`
//	// APIKey токен для работы с проектом
//	APIKey string `json:"api_key,omitempty"`
//	// ProjectName название проекта
//	ProjectName string `json:"project_name"`
//	Jobs        Jobs   `json:"jobs,omitempty"`
//}

type Response struct {
	Message string  `json:"message"`
	Error   *string `json:"error,omitempty"`
}

type ProjectResponse struct {
	Project *Project `json:"project"`
	Message string   `json:"message,omitempty"`
	Error   *string  `json:"error,omitempty"`
}

type ProjectsResponse struct {
	Projects []*Project `json:"projects"`
	Message  string     `json:"message,omitempty"`
	Error    *string    `json:"error,omitempty"`
}

type JobResponse struct {
	Job     *Job    `json:"job"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

type JobsResponse struct {
	Jobs    []*Job  `json:"jobs"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}
