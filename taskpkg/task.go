package taskpkg

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/sources"
	"path/filepath"
	"strings"
)

//var (
//	taskWorkerInstance TaskWorker
//)
//
//func GetWorker() *TaskWorker {
//	return &taskWorkerInstance
//}
//
//func NewWorker(hub worker.IWorkerHub) *TaskWorker {
//	taskWorkerInstance = TaskWorker{
//		TaskType:    CWorkerType,
//		Running:     false,
//		controlChan: make(chan int, 1),
//		hub:         hub,
//	}
//	return &taskWorkerInstance
//}

func PrepareStepCommand(j sources.Job, s pipeline.Step) (string, error) {
	var commands []string

	// Базовые команды
	commands = append(commands, "mkdir /job")
	commands = append(commands, "cd /job")
	commands = append(commands, "git clone "+j.URL)

	// Извлекаем путь из URL
	repoPath := strings.TrimSuffix(j.URL, ".git") // Удаляем расширение .git, если оно есть
	repoName := filepath.Base(repoPath)

	commands = append(commands, "cd "+repoName)

	if j.Branch != "" {
		commands = append(commands, "git switch "+j.Branch)
	}

	// Команды из пайплайна
	for _, c := range s.Commands {
		commands = append(commands, c)
	}

	return strings.Join(commands, " && "), nil
}

// PrepareSubTasksStruct подготовить структуру шага
func PrepareSubTasksStruct(s pipeline.Pipeline) []SubTaskResult {
	var subTasks []SubTaskResult
	for _, step := range s.Steps {
		subTasks = append(subTasks, SubTaskResult{
			Status: StepPending,
			Name:   step.Name,
		})
	}

	return subTasks
}
