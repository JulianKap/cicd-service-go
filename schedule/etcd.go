package schedule

import (
	"cicd-service-go/db/etcd"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// getTasksETCD получить список всех тасок из etcd
func getTasksETCD(cli *clientv3.Client, p *sources.Project, tasks *taskpkg.Tasks) error {
	resp, err := etcd.GetKey(cli, Keys.Tasks)
	if err != nil {
		log.Error("getTasksETCD #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getTasksETCD #1: key ", Keys.Tasks, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getTasksETCD #2: key ", Keys.Tasks, " get more than one key")
	}

	var allTask taskpkg.Tasks
	if err := json.Unmarshal(resp.Kvs[0].Value, &allTask); err != nil {
		log.Error("getTasksETCD #3: ", err)
		return err
	}

	for _, t := range allTask.Tasks {
		if t.ProjectID == p.ID {
			tasks.Tasks = append(tasks.Tasks, t)
		}
	}

	return nil
}

// getTaskETCD
func getTaskETCD(cli *clientv3.Client, p *sources.Project, task *taskpkg.Task) (bool, error) {
	return true, nil
}

// deleteTaskETCD
func deleteTaskETCD(cli *clientv3.Client, p *sources.Project, task *taskpkg.Task) (bool, error) {
	return true, nil
}
