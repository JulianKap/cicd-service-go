package pipeline

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

	if branch == "" {
		branch = "master"
	}

	if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
	}
	if strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	}

	// Формируем URL для получения cicd.yml файла
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/cicd.yml?ref=%s", url, branch)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("getPipelineFromGit #3: ", err)
		return nil, err
	}

	var pRaw PipelineRawString
	err = json.Unmarshal(body, &pRaw)
	if err != nil {
		log.Error("getPipelineFromGit #4: ", err)
		return nil, err
	}

	// Декодирование содержимого файла из base64
	content, err := base64.StdEncoding.DecodeString(pRaw.Content)
	if err != nil {
		log.Error("getPipelineFromGit #5: ", err)
		return nil, err
	}

	return content, nil
}
