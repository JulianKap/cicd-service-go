package sources

// Project представляет информацию о проекте
type Project struct {
	ID          string `json:"id"`
	APIKey      string `json:"api_key"`
	ProjectName string `json:"project_name"`
	//Jobs   map[string]Job `json:"jobs"`
}

// Job представляет информацию о задании
type Job struct {
	ID        string `json:"id"`
	JobName   string `json:"jobName"` // Название задачи
	IdProject string `json:"idProject"`
	URL       string `json:"url"` // URL до репозитория
	//Credentials string `json:"credentials"`	// Креды для доступа к репозиторию todo: отдельную структуру для кредов
	Branch string `json:"branch"` // Ветка репозитория
}

type Response struct {
	Message string `json:"message,omitempty"`
}

//todo: добавить валидации на ключевые поля в структурах
