package taskpkg

// Task представляет информацию о текущей задаче
type Task struct {
	ID        string `json:"id"`         // Уникальный идентификатор задачи
	ProjectID string `json:"project_id"` // Идентификатор проекта, к которому относится задача
	JobID     string `json:"job_id"`     // Идентификатор задания, выполняемого задачей
	Status    string `json:"status"`     // Статус задачи (например, "running", "finished", "failed" и т.д.)
}
