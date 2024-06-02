package taskpkg

import (
	"cicd-service-go/pipeline"
	"cicd-service-go/sources"
	"path/filepath"
	"strings"
)

func PrepareStepCommand(j sources.Job, s pipeline.Step) (string, error) {
	var commands []string

	// Базовые команды
	commands = append(commands, "mkdir /job")
	commands = append(commands, "cd /job")

	if j.Branch != "" {
		commands = append(commands, "git clone --branch "+j.Branch+" "+j.URL)
	} else {
		commands = append(commands, "git clone "+j.URL)
	}

	// Извлекаем путь из URL
	repoPath := strings.TrimSuffix(j.URL, ".git") // Удаляем расширение .git, если оно есть
	repoName := filepath.Base(repoPath)

	commands = append(commands, "cd "+repoName)

	// Команды из пайплайна
	for _, c := range s.Commands {
		commands = append(commands, c)
	}

	result := ""
	for _, c := range commands {
		result += c + " || exit 1\n"
	}

	//return strings.Join(commands, " && "), nil
	return result, nil
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
