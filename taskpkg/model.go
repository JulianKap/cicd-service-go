package taskpkg

import (
	"time"
)

// KeysDCS ключи в DCS
type KeysDCS struct {
	// Tasks все таски
	Tasks string `json:"tasks"`
	// TasksHistory история выполнения тасок
	TasksHistory string `json:"tasks_history"`
	// TaskProject таски конкретного проекта
	TaskProject string `json:"task_project"`
	// TaskLatestId крайний id таски
	TaskLatestId string `json:"task_latest_id"`
}

// TaskStatus представляет статус выполнения задания
type TaskStatus int

const (
	Pending   TaskStatus = iota // Задание ожидает выполнения
	Running                     // Задание выполняется
	Completed                   // Задание завершена успешно
	Failed                      // Задание завершена с ошибкой
	Schedule                    // Задание распределено, но еще не запускалось
)

// TaskResult содержит результат выполнения задания
type TaskResult struct {
	// Status статус выполнения задачи
	Status TaskStatus `json:"status"`
	// Message сообщение о результате выполнения задачи
	Message string `json:"message"`
	// RetryCount количество попыток повторного запуска
	RetryCount int `json:"retry_count"`
	// CreateAt время создания таски
	RunningAt *time.Time `json:"running_at"`
	// WorkerUUID uuid воркера, на которого распределено задание
	WorkerUUID string `json:"worker_uuid"`
}

// Task представляет информацию о текущей таске
type Task struct {
	// ID уникальный идентификатор таски
	ID int `json:"id"`
	// ProjectID идентификатор проекта, к которому относится таска
	ProjectID int `json:"project_id"`
	// JobID идентификатор задания, которое запускает таска
	JobID int `json:"job_id"`
	// Name название таски
	Name string `json:"name,omitempty"`
	// Status статус выполнения (например, "running", "finished", "failed" и т.д.)
	Status TaskResult `json:"status"`
	// CreateAt время создания таски
	CreateAt *time.Time `json:"create_at"`
	// NumberOfRetriesOnError количество попыток для повторного запуска при ошибке
	NumberOfRetriesOnError int `json:"number_of_retries_on_error"`
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

//
//type ITask interface {
//	GetActive() error
//	IsZeroID() bool
//	SetStatus(string)
//	Save() error
//	Finish() error
//	GetJSON() ([]byte, error)
//}
//
//type TaskWorker struct {
//	TaskType    string
//	TaskName    string
//	Running     bool
//	controlChan chan int `json:"-"`
//
//	//hub         worker.IWorkerHub
//	CurrentTask *Task
//
//	cancel context.CancelFunc
//}
