package pipeline

import (
	"cicd-service-go/sources"
	log "github.com/sirupsen/logrus"
)

func GetPipeline(j sources.Job) (p Pipeline, err error) {
	// Получаем содержимое файла cicd.yml из репозитория git
	content, err := getPipelineFromGit(j.URL, j.Branch)
	if err != nil {
		log.Error("GetPipeline #0: ", err)
		return p, err
	}

	// Парсим содержимое файла в структуру CICDPipeline
	p, err = parsePipeline(content)
	if err != nil {
		log.Error("GetPipeline #1: ", err)
		return p, err
	}

	return p, nil
}
