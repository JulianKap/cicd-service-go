package taskpkg

import (
	"time"
)

// KeysDCS ключи в DCS
type KeysDCS struct {
	// Tasks все таски
	Tasks string `json:"tasks"`
	// TasksHistory история выполения тасок
	TasksHistory string `json:"tasks_history"`
	// TaskProject таски конкретного проекта
	TaskProject string `json:"task_project"`
	// TaskLatestId крайний id таски
	TaskLatestId string `json:"task_latest_id"`
}

// Task представляет информацию о текущей таске
type Task struct {
	// ID уникальный идентификатор задачи
	ID int `json:"id"`
	// ProjectID идентификатор проекта, к которому относится задача
	ProjectID int `json:"project_id"`
	// JobID идентификатор задания, выполняемого задачей
	JobID int `json:"job_id"`
	// Name название задачи
	Name string `json:"name,omitempty"`
	// Status статус задачи (например, "running", "finished", "failed" и т.д.)
	Status string `json:"status"`
	// AddAt время создания задачи
	CreateAt *time.Time `json:"create_at"`
}

// Tasks список тасок
type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type TaskResponse struct {
	Task    *Task   `json:"task"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks   *Tasks  `json:"tasks"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}
