package schedule

import (
	"cicd-service-go/db/etcd"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"encoding/json"
	"strconv"

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

// getTaskETCD получение таски по id
func getTaskETCD(cli *clientv3.Client, t *taskpkg.Task) (bool, error) {
	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTaskETCD #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getTaskETCD #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getTaskETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &t); err != nil {
		log.Error("getTaskETCD #3: ", err)
		return false, err
	}

	return true, nil
}

// createTaskETCD создание таски
func createTaskETCD(cli *clientv3.Client, p *sources.Project, t *taskpkg.Task) error {
	latestId, err := sources.GetLatestIdETCD(cli, Keys.TaskLatestId)
	if err != nil {
		log.Error("createTaskETCD #0: ", err)
		return err
	}

	if latestId == -1 {
		latestId = 1
	}
	t.ID = latestId + 1

	// Добавление тасок в общий список
	var tasks taskpkg.Tasks
	if err := getTasksETCD(cli, p, &tasks); err != nil {
		log.Error("createTaskETCD #1: ", err)
		return err
	}
	tasks.Tasks = append(tasks.Tasks, *t)

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Error("createTaskETCD #2: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Tasks, string(tasksJSON)); err != nil {
		log.Error("createTaskETCD #3: ", err)
		return err
	}

	// Добавление таски к своему проекту
	taskJSON, err := json.Marshal(t)
	if err != nil {
		log.Error("createTaskETCD #4: ", err)
		return err
	}

	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	if err = etcd.SetKey(cli, key, string(taskJSON)); err != nil {
		log.Error("createTaskETCD #5: ", err)
		return err
	}

	// Добавление последнего ID
	if err = etcd.SetKey(cli, Keys.TaskLatestId, strconv.Itoa(t.ID)); err != nil {
		log.Error("createTaskETCD #6: ", err)
		return err
	}

	return nil
}

// deleteTaskETCD удаление таски
func deleteTaskETCD(cli *clientv3.Client, p *sources.Project, t *taskpkg.Task) (bool, error) {
	var tasks taskpkg.Tasks
	if err := getTasksETCD(cli, p, &tasks); err != nil {
		log.Error("deleteTaskETCD #0: ", err)
		return false, err
	}

	state := false
	var newTasks taskpkg.Tasks
	for _, task := range tasks.Tasks {
		if t.ID != task.ID {
			newTasks.Tasks = append(newTasks.Tasks, task)
		} else {
			state = true
		}
	}

	valueJSON, err := json.Marshal(newTasks)
	if err != nil {
		log.Error("deleteTaskETCD #1: ", err)
		return state, err
	}

	// Обновляем список всех тасок
	if err = etcd.SetKey(cli, Keys.Tasks, string(valueJSON)); err != nil {
		log.Error("deleteTaskETCD #2: ", err)
		return state, err
	}

	// Удаляем таску
	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	if err = etcd.DelKey(cli, key); err != nil {
		log.Error("deleteTaskETCD #3: ", err)
		return state, err
	}

	return state, nil
}
