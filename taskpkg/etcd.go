package taskpkg

import (
	"cicd-service-go/constants"
	"cicd-service-go/db/etcd"
	"cicd-service-go/manager"
	"cicd-service-go/sources"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
)

// GetKeyTasks получение ключ всех заданий для всех проектов
func GetKeyTasks() string {
	return manager.Conf.Cluster.Namespace + constants.PROJECTS_TASKS
}

// GetKeyTasksHistory получение ключ истории всех заданий
func GetKeyTasksHistory() string {
	return manager.Conf.Cluster.Namespace + constants.PROJECTS_TASKS_HISTORY
}

// GetKeyLatestIdTask получить ключ идентификатора последнего задания
func GetKeyLatestIdTask() string {
	return manager.Conf.Cluster.Namespace + constants.PROJECTS_TASKS_LATEST_ID
}

// GetKeyTaskProject получение конкретного задания у конкретного проекта
func GetKeyTaskProject(t *Task) string {
	return sources.GetKeyProject(&sources.Project{ID: t.ProjectID}) + constants.TASKS + "/" + strconv.Itoa(t.ID)
}

// getTasksETCD получить список всех заданий из etcd по всем проектам
func (t *Tasks) getTasksETCD(cli *clientv3.Client, isHistory bool) error {
	var key string
	if isHistory {
		key = GetKeyTasksHistory()
	} else {
		key = GetKeyTasks()
	}

	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTasksETCD #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getTasksETCD #1: key ", key, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getTasksETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &t); err != nil {
		log.Error("getTasksETCD #3: ", err)
		return err
	}

	return nil
}

// setTasksETCD добавить список заданий для всех проектов
func (t *Tasks) setTasksETCD(cli *clientv3.Client, isHistory bool) error {
	if len(t.Tasks) == 0 {
		log.Info("setTasksETCD #0: not found tasks")
	}

	tasksJSON, err := json.Marshal(t)
	if err != nil {
		log.Error("setTasksETCD #1: ", err)
		return err
	}

	var key string
	if isHistory {
		key = GetKeyTasksHistory()
	} else {
		key = GetKeyTasks()
	}

	if err = etcd.SetKey(cli, key, string(tasksJSON)); err != nil {
		log.Error("setTasksETCD #3: ", err)
		return err
	}

	return nil
}

// getTasksByProjectETCD получить список всех заданий из etcd по указанному проекту
func (t *Tasks) getTasksByProjectETCD(cli *clientv3.Client, p *sources.Project, isHistory bool) error {
	var allTask Tasks
	if err := allTask.getTasksETCD(cli, isHistory); err != nil {
		log.Error("getTasksByProjectETCD #0: ", err)
		return err
	}

	for _, task := range allTask.Tasks {
		if task.ProjectID == p.ID {
			t.Tasks = append(t.Tasks, task)
		}
	}

	return nil
}

// getTaskByProjectETCD получить указанное задание по проекту
func (t *Task) getTaskByProjectETCD(cli *clientv3.Client) (bool, error) {
	key := GetKeyTaskProject(t)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTaskByProjectETCD #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getTaskByProjectETCD #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getTaskByProjectETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &t); err != nil {
		log.Error("getTaskByProjectETCD #3: ", err)
		return false, err
	}

	return true, nil
}

// setTaskByProjectETCD создание задания
func (t *Task) setTaskByProjectETCD(cli *clientv3.Client, p *sources.Project) error {
	latestId, err := etcd.GetKeyInt(cli, GetKeyLatestIdTask())
	if err != nil {
		log.Error("setTaskByProjectETCD #0: ", err)
		return err
	}

	if latestId == -1 {
		latestId = 1
	}
	t.ID = latestId + 1

	// Добавление тасок в общий список
	var tasks Tasks
	if err := tasks.getTasksByProjectETCD(cli, p, false); err != nil {
		log.Error("setTaskByProjectETCD #1: ", err)
		return err
	}
	// Добавление задания в общий список заданий
	tasks.Tasks = append(tasks.Tasks, t)
	if err = tasks.setTasksETCD(cli, false); err != nil {
		log.Error("setTaskByProjectETCD #2: ", err)
		return err
	}

	// Добавление задания к своему проекту
	taskJSON, err := json.Marshal(t)
	if err != nil {
		log.Error("setTaskByProjectETCD #3: ", err)
		return err
	}

	if err = etcd.SetKey(cli, GetKeyTaskProject(t), string(taskJSON)); err != nil {
		log.Error("setTaskByProjectETCD #4: ", err)
		return err
	}

	// Добавление последнего ID задания
	if err = etcd.SetKey(cli, GetKeyLatestIdTask(), strconv.Itoa(t.ID)); err != nil {
		log.Error("setTaskByProjectETCD #5: ", err)
		return err
	}

	return nil
}

// markerDelTaskByProjectETCD отметить задачу как удаляемую. Планировщик все проверит и сам ее удалит
func (t *Task) markerDelTaskByProjectETCD(cli *clientv3.Client, p *sources.Project) (bool, error) {
	var tasks Tasks
	if err := tasks.getTasksByProjectETCD(cli, p, false); err != nil {
		log.Error("markerDelTaskByProjectETCD #0: ", err)
		return false, err
	}

	state := false
	for _, task := range tasks.Tasks {
		if t.ID == task.ID {
			task.Status.Status = Removing //Отмечаем, что задание нужно удалить
			state = true
		}
	}

	valueJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Error("markerDelTaskByProjectETCD #1: ", err)
		return state, err
	}

	// Обновляем список всех тасок
	if err = etcd.SetKey(cli, GetKeyTasks(), string(valueJSON)); err != nil {
		log.Error("markerDelTaskByProjectETCD #2: ", err)
		return state, err
	}

	//// Удаляем таску
	//if err = etcd.DelKey(cli, GetKeyTaskProject(t)); err != nil {
	//	log.Error("markerDelTaskByProjectETCD #3: ", err)
	//	return state, err
	//}

	return state, nil
}

// setHistoryTaskByProjectETCD добавление задания в список истории
func (t *Task) setHistoryTaskByProjectETCD(cli *clientv3.Client, p *sources.Project) error {
	var historyTasks Tasks
	if err := historyTasks.getTasksByProjectETCD(cli, p, true); err != nil {
		log.Error("setHistoryTaskByProjectETCD #1: ", err)
		return err
	}
	// Добавление задания в список истории
	historyTasks.Tasks = append(historyTasks.Tasks, t)
	if err := historyTasks.setTasksETCD(cli, true); err != nil {
		log.Error("setHistoryTaskByProjectETCD #2: ", err)
		return err
	}

	return nil
}

func GetActualTasksByProjectETCD(cli *clientv3.Client, p *sources.Project, t *Tasks) error {
	return t.getTasksByProjectETCD(cli, p, false)
}

func SetActualTasksETCD(cli *clientv3.Client, t *Tasks) error {
	return t.setTasksETCD(cli, false)
}

func GetActualTasksETCD(cli *clientv3.Client, t *Tasks) error {
	return t.getTasksETCD(cli, false)
}

func SetHistoryTaskByProjectETCD(cli *clientv3.Client, p *sources.Project, t *Task) error {
	return t.setHistoryTaskByProjectETCD(cli, p)
}
