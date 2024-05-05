package schedule

import (
	"cicd-service-go/constants"
	"cicd-service-go/db/etcd"
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// setTasksETCD добавить список заданий для всех проектов
func setTasksETCD(cli *clientv3.Client, tasks *taskpkg.Tasks) error {
	if len(tasks.Tasks) == 0 {
		log.Info("setTasksETCD #0: not found tasks")
	}

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Error("setTasksETCD #1: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Tasks, string(tasksJSON)); err != nil {
		log.Error("setTasksETCD #3: ", err)
		return err
	}

	return nil
}

// getTasksETCD получить список всех заданий из etcd по всем проектам
func getTasksETCD(cli *clientv3.Client, tasks *taskpkg.Tasks) error {
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

	if err := json.Unmarshal(resp.Kvs[0].Value, &tasks); err != nil {
		log.Error("getTasksETCD #3: ", err)
		return err
	}

	return nil
}

// getTasksByProjectETCD получить список всех заданий из etcd по указанному проекту
func getTasksByProjectETCD(cli *clientv3.Client, p *sources.Project, tasks *taskpkg.Tasks) error {
	var allTask taskpkg.Tasks
	if err := getTasksETCD(cli, &allTask); err != nil {
		log.Error("getTasksByProjectETCD #0: ", err)
		return err
	}

	for _, t := range allTask.Tasks {
		if t.ProjectID == p.ID {
			tasks.Tasks = append(tasks.Tasks, t)
		}
	}

	return nil
}

// getTaskByProjectETCD получить указанное задание по проекту
func getTaskByProjectETCD(cli *clientv3.Client, t *taskpkg.Task) (bool, error) {
	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTaskByProjectETCD #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getTaskByProjectETCD #1: key ", key, " not found")
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

// createTaskByProjectETCD создание задания
func createTaskByProjectETCD(cli *clientv3.Client, p *sources.Project, t *taskpkg.Task) error {
	latestId, err := sources.GetLatestIdETCD(cli, Keys.TaskLatestId)
	if err != nil {
		log.Error("createTaskByProjectETCD #0: ", err)
		return err
	}

	if latestId == -1 {
		latestId = 1
	}
	t.ID = latestId + 1

	// Добавление тасок в общий список
	var tasks taskpkg.Tasks
	if err := getTasksByProjectETCD(cli, p, &tasks); err != nil {
		log.Error("createTaskByProjectETCD #1: ", err)
		return err
	}
	// Добавление задания в общий список заданий
	tasks.Tasks = append(tasks.Tasks, *t)
	if err = setTasksETCD(cli, &tasks); err != nil {
		log.Error("createTaskByProjectETCD #2: ", err)
		return err
	}

	// Добавление задания к своему проекту
	taskJSON, err := json.Marshal(t)
	if err != nil {
		log.Error("createTaskByProjectETCD #3: ", err)
		return err
	}

	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	if err = etcd.SetKey(cli, key, string(taskJSON)); err != nil {
		log.Error("createTaskByProjectETCD #4: ", err)
		return err
	}

	// Добавление последнего ID задания
	if err = etcd.SetKey(cli, Keys.TaskLatestId, strconv.Itoa(t.ID)); err != nil {
		log.Error("createTaskByProjectETCD #5: ", err)
		return err
	}

	return nil
}

// deleteTaskByProjectETCD удаление задания
func deleteTaskByProjectETCD(cli *clientv3.Client, p *sources.Project, t *taskpkg.Task) (bool, error) {
	var tasks taskpkg.Tasks
	if err := getTasksByProjectETCD(cli, p, &tasks); err != nil {
		log.Error("deleteTaskByProjectETCD #0: ", err)
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
		log.Error("deleteTaskByProjectETCD #1: ", err)
		return state, err
	}

	// Обновляем список всех тасок
	if err = etcd.SetKey(cli, Keys.Tasks, string(valueJSON)); err != nil {
		log.Error("deleteTaskByProjectETCD #2: ", err)
		return state, err
	}

	// Удаляем таску
	key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
	if err = etcd.DelKey(cli, key); err != nil {
		log.Error("deleteTaskByProjectETCD #3: ", err)
		return state, err
	}

	return state, nil
}

// getTasksForWorker получить список всех заданий для воркера
func getTasksForWorker(cli *clientv3.Client, m manager.Member, t *taskpkg.Tasks) (err error) {
	key := manager.Conf.Cluster.Namespace + constants.WORKERS + "/" + m.UUID + "/tasks"
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTasksForWorker #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getTasksForWorker #1: key ", key, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getTasksForWorker #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &t); err != nil {
		log.Error("getTasksForWorker #3: ", err)
		return err
	}

	return nil
}

// getTaskForWorker получение задания воркера с актуальным статусом его выполнения
func getTaskForWorker(cli *clientv3.Client, m manager.Member, t *taskpkg.Task) error {
	var tasks taskpkg.Tasks
	if err := getTasksForWorker(cli, m, &tasks); err != nil {
		log.Error("getTaskForWorker #0: ", err)
		return err
	}

	for _, task := range tasks.Tasks {
		if task.ID == t.ID {
			t.Status = task.Status
			break
		}
	}

	return nil
}

// setTaskForWorker назначить задание для воркера
func setTaskForWorker(cli *clientv3.Client, m manager.Member, t *taskpkg.Task) (state bool, err error) {
	var tasks taskpkg.Tasks
	if err := getTasksForWorker(cli, m, &tasks); err != nil {
		log.Error("setTaskToWorker #0: ", err)
		return false, err
	}

	for _, task := range tasks.Tasks {
		if task.ID == t.ID {
			log.Info("setTaskToWorker #1: task id: ", t.ID, " already exists")
			return true, nil
		}
	}

	tasks.Tasks = append(tasks.Tasks, *t)

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Error("setTaskToWorker #2: ", err)
		return false, err
	}

	key := manager.Conf.Cluster.Namespace + constants.WORKERS + "/" + m.UUID + "/tasks"
	if err = etcd.SetKey(cli, key, string(tasksJSON)); err != nil {
		log.Error("setTaskToWorker #3: ", err)
		return false, err
	}

	return true, nil
}

// updateTaskForWorker обновить статус выполнения задания для воркера
func updateTaskForWorker(cli *clientv3.Client, m manager.Member, t *taskpkg.Task) error {
	var tasks taskpkg.Tasks
	if err := getTasksForWorker(cli, m, &tasks); err != nil {
		log.Error("updateTaskForWorker #0: ", err)
		return err
	}

	var newTasks taskpkg.Tasks
	for _, task := range tasks.Tasks {
		if task.ID == t.ID {
			log.Debug("updateTaskForWorker #1: update status task on ", t.Status.Status, " (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
			newTasks.Tasks = append(newTasks.Tasks, *t)
		} else {
			newTasks.Tasks = append(newTasks.Tasks, task)
		}
	}

	tasksJSON, err := json.Marshal(newTasks)
	if err != nil {
		log.Error("updateTaskForWorker #2: ", err)
		return err
	}

	key := manager.Keys.Worker + "/tasks"
	if err = etcd.SetKey(cli, key, string(tasksJSON)); err != nil {
		log.Error("updateTaskForWorker #3: ", err)
		return err
	}

	return nil
}

// updateAllTasks обновление всех заданий
func updateAllTasks(cli *clientv3.Client, tasks *taskpkg.Tasks) error {
	// Обновление списка заданий по всем проектам
	if err := setTasksETCD(cli, tasks); err != nil {
		log.Error("updateAllTasks #0: ", err)
		return err
	}

	// Обновление заданий в рамках проекта
	for _, t := range tasks.Tasks {
		// Добавление таски к своему проекту
		taskJSON, err := json.Marshal(t)
		if err != nil {
			log.Error("updateAllTasks #1: ", err)
			continue
		}

		key := Keys.TaskProject + "/" + strconv.Itoa(t.ProjectID) + "/tasks/" + strconv.Itoa(t.ID)
		if err = etcd.SetKey(cli, key, string(taskJSON)); err != nil {
			log.Error("updateAllTasks #2: ", err)
			continue
		}
	}

	return nil
}

// getJobEtcd получить задачу, для которой относится задание
func getJobEtcd(t taskpkg.Task) (job sources.Job, err error) {
	p := sources.Project{ID: t.ProjectID}
	j := sources.Job{ID: t.JobID}

	if _, err := p.GetJobETCD(db.InstanceETCD, &j); err != nil {
		log.Error("getJobEtcd #0: ", err)
		return j, err
	}

	return j, nil
}
