package taskpkg

import "time"

// Task представляет информацию о текущей задаче
type Task struct {
	// ID уникальный идентификатор задачи
	ID string `json:"id"`
	// ProjectID идентификатор проекта, к которому относится задача
	ProjectID string `json:"project_id"`
	// JobID идентификатор задания, выполняемого задачей
	JobID string `json:"job_id"`
	// Name название задачи
	Name string `json:"name,omitempty"`
	// Status статус задачи (например, "running", "finished", "failed" и т.д.)
	Status string `json:"status"`
	// AddAt время создания задачи
	CreateAt *time.Time `json:"create_at"`
}
