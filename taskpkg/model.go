package taskpkg

import (
	"time"
)

// TaskStatus представляет статус выполнения задания
type TaskStatus int
type SubTaskStatus int

const (
	Pending   TaskStatus = iota // Задание ожидает выполнения
	Running                     // Задание выполняется
	Completed                   // Задание завершена успешно
	Failed                      // Задание завершена с ошибкой
	Schedule                    // Задание распределено, но еще не запускалось
	Removing                    // Задание, которое нужно удалить
)

const (
	StepSkipped   SubTaskStatus = iota // Шаг пропущен
	StepRunning                        // Шаг запущен
	StepComplited                      // Шаг выполнен
	StepFailed                         // Шаг с ошибкой
	StepPending                        // Шаг в ожидании запуска
)

// SubTaskResult содержит результат выполнения шага задания (step в пайплайне)
type SubTaskResult struct {
	// Status статус выполнения шага
	Status SubTaskStatus `json:"status"`
	// Message сообщение о результате выполнения задания
	Message string `json:"message,omitempty"`
	// FinnishAt время выполнения шага
	ElapsedSec int `json:"elapsed_sec,omitempty"`
	// Name название шага
	Name string `json:"name"`
}

// TaskResult содержит результат выполнения задания
type TaskResult struct {
	// Status статус выполнения задания
	Status TaskStatus `json:"status"`
	// Message сообщение о результате выполнения задания
	Message string `json:"message,omitempty"`
	// RetryCount количество запусков задания
	RetryCount int `json:"retry_count"`
	// CreateAt время запуска задания
	RunningAt *time.Time `json:"running_at,omitempty"`
	// FinnishAt время выполнения задания
	ElapsedSec int `json:"elapsed_sec,omitempty"`
	// WorkerUUID uuid воркера, на которого распределено задание
	WorkerUUID string `json:"worker_uuid"`
	// Steps шаги задания
	Steps []SubTaskResult `json:"steps,omitempty"`
}

// Task представляет информацию о текущем задании
type Task struct {
	// ID уникальный идентификатор задания
	ID int `json:"id"`
	// ProjectID идентификатор проекта, к которому относится задание
	ProjectID int `json:"project_id"`
	// JobID идентификатор задачи, которое запускает задание
	JobID int `json:"job_id"`
	// Name название задания
	Name string `json:"name,omitempty"`
	// Status статус выполнения задания
	Status TaskResult `json:"status,omitempty"`
	// CreateAt время создания задания
	CreateAt *time.Time `json:"create_at,omitempty"`
	// NumberOfRetriesOnError количество попыток для повторного запуска при ошибке
	NumberOfRetriesOnError int `json:"number_of_retries_on_error"`
}

// Tasks список заданий
type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

type TaskResponse struct {
	Task    *Task   `json:"task"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks   []*Task `json:"tasks"`
	Message string  `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}
