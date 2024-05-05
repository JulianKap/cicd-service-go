package pipeline

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// getPipelineFromGit получение пайплайна в указанном репозитории репозитория
func getPipelineFromGit(url string, branch string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	}
	// Формируем URL для получения cicd.yml файла
	apiURL := fmt.Sprintf("%s/%s/cicd.yml", url, branch)

	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		log.Error("getPipelineFromGit #0: ", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("getPipelineFromGit #1: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("getPipelineFromGit #2: ", err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("getPipelineFromGit #3: ", err)
		return nil, err
	}

	return data, nil
}
