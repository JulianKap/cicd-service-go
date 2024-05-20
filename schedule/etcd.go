package schedule

import (
	"cicd-service-go/constants"
	"cicd-service-go/db/etcd"
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// getTasksForWorker получить список всех заданий для воркера
func getTasksForWorker(cli *clientv3.Client, m manager.Member, t *taskpkg.Tasks) (err error) {
	key := manager.Conf.Cluster.Namespace + constants.WORKERS + "/" + m.UUID + "/tasks" //todo: править
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getTasksForWorker #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getTasksForWorker #1: key ", key, " not found")
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

	tasks.Tasks = append(tasks.Tasks, t)

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Error("setTaskToWorker #2: ", err)
		return false, err
	}

	key := manager.Conf.Cluster.Namespace + constants.WORKERS + "/" + m.UUID + "/tasks" //todo: править
	if err = etcd.SetKey(cli, key, string(tasksJSON)); err != nil {
		log.Error("setTaskToWorker #3: ", err)
		return false, err
	}

	return true, nil
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
			newTasks.Tasks = append(newTasks.Tasks, t)
		} else {
			newTasks.Tasks = append(newTasks.Tasks, task)
		}
	}

	tasksJSON, err := json.Marshal(newTasks)
	if err != nil {
		log.Error("updateTaskForWorker #2: ", err)
		return err
	}

	key := manager.Keys.Worker + "/tasks" //todo: править
	if err = etcd.SetKey(cli, key, string(tasksJSON)); err != nil {
		log.Error("updateTaskForWorker #3: ", err)
		return err
	}

	return nil
}

// updateAllTasks обновление всех заданий
func updateAllTasks(cli *clientv3.Client, tasks *taskpkg.Tasks) error {
	// Обновление списка заданий по всем проектам
	if err := taskpkg.SetTasksETCD(cli, tasks); err != nil {
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

		if err = etcd.SetKey(cli, taskpkg.GetKeyTaskProject(t), string(taskJSON)); err != nil {
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
