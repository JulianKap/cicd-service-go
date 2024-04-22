package schedule

import (
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/taskpkg"
	log "github.com/sirupsen/logrus"
)

// tasksScheduler обработка тасок мастером
func tasksScheduler() (standalone bool, err error) {
	standalone = false

	var tasks taskpkg.Tasks
	if err := getTasksETCD(db.InstanceETCD, &tasks); err != nil {
		log.Error("tasksScheduler #0: ", err)
		return standalone, err
	}

	if len(tasks.Tasks) == 0 {
		log.Debug("tasksScheduler #1: not found tasks")
		return standalone, nil
	}

	members, err := manager.GetMembers()
	if err != nil {
		log.Error("tasksScheduler #2: ", err)
		return standalone, err
	}

	if len(members.Members) == 0 {
		log.Warn("tasksScheduler #3: not found members")
		return standalone, nil
	} else if len(members.Members) == 1 {
		if members.Members[0].UUID != manager.MemberInfo.UUID {
			log.Warn("tasksScheduler #4: one member in the cluster and is not me")
			return standalone, nil
		}
		standalone = true

		log.Debug("tasksScheduler #5: one member in cluster")
	}

	tasksInQueue := taskpkg.Tasks{}
	for _, t := range tasks.Tasks {
		if t.Status.Status == taskpkg.Completed {
			// Перенести в history
			if err := moveTaskInHistory(t); err != nil {
				log.Error("tasksScheduler #6: ", err)
			}
		} else if t.Status.Status == taskpkg.Failed {
			// Лимит попыток исчерпан. Переводим задачу в список истории
			if t.NumberOfRetriesOnError > t.Status.RetryCount {
				tasksInQueue.Tasks = append(tasksInQueue.Tasks, t)
				continue
			}

			if err := moveTaskInHistory(t); err != nil {
				log.Error("tasksScheduler #6: ", err)
			}
		} else if t.Status.Status == taskpkg.Running {
			// Проверить статус выполнения
			// todo: проверить, что задача действительно выполняется, если же нет, то нужно либо добавить в очередь выполнения
			// либо перевести в историю. Данный пункт нужно особенно тщательно протестировать.
			// К примеру добавить дедлайн на выполнения, и если превышение, то отменяем задачу.

		} else if t.Status.Status == taskpkg.Pending {
			tasksInQueue.Tasks = append(tasksInQueue.Tasks, t)
		}
	}

	if len(tasksInQueue.Tasks) == 0 {
		log.Debug("tasksScheduler #7: no tasks to run")
		return true, nil
	}

	if standalone {
		for _, t := range tasksInQueue.Tasks {
			ok, err := setTaskWorker(members.Members[0], t)
			if err != nil {
				log.Error("tasksScheduler #8: ", err)
				//return standalone, nil
			}

			if !ok {
				// todo: как вариат, отмечать или просто пропускать такие таски
			}
		}
	} else {
		i := 0
		// Распределяем задачи по воркерам
		for _, t := range tasksInQueue.Tasks {
			ok, err := setTaskWorker(members.Members[i], t)
			if err != nil {
				log.Error("tasksScheduler #9: ", err)
			}

			if !ok {
				// todo: распределить данную таску на другого воркера
			}

			i++
			if i == len(members.Members) {
				i = 0
			}
		}
	}

	return true, err
}

// setTaskWorker назначение задачи члену кластера
func setTaskWorker(w manager.Member, t taskpkg.Task) (ok bool, err error) {
	log.Debug("setTaskWorker #0: member: ", w.UUID, " task: ", t.ID)
	return setTaskToWorker(db.InstanceETCD, w, &t)
}

// moveTaskInHistory перевод таски в список истории
func moveTaskInHistory(task taskpkg.Task) error {
	log.Debug("setTaskWorker #0: task: ", task.ID)

	// todo: добавить отправку в task_history
	return nil
}

// tasksWorker обработка задач воркером
func tasksWorker() (err error) {

	// Получаем список задач

	// Проверяем, если

	//

	return err
}
