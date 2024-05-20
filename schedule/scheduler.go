package schedule

import (
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	log "github.com/sirupsen/logrus"
)

// tasksScheduler планирование всех заданий мастером
func tasksScheduler() (bool, error) {
	standalone := false

	// Получаем список всех тасок
	var tasks taskpkg.Tasks
	if err := taskpkg.GetTasksETCD(db.InstanceETCD, &tasks); err != nil {
		log.Error("tasksScheduler #0: ", err)
		return false, err
	}

	if len(tasks.Tasks) == 0 {
		log.Debug("tasksScheduler #1: not found tasks")
		return false, nil
	}

	// Получаем список всех членов кластера
	members, err := manager.GetMembers()
	if err != nil {
		log.Error("tasksScheduler #2: ", err)
		return false, err
	}

	workers := manager.Members{}

	if len(members.Members) == 0 {
		log.Warn("tasksScheduler #3: not found members")
		return false, nil
	} else if len(members.Members) == 1 {
		if members.Members[0].UUID != manager.MemberInfo.UUID {
			log.Warn("tasksScheduler #4: one member in the cluster and is not me")
			return false, nil
		}
		standalone = true

		log.Debug("tasksScheduler #5: one member in cluster")
	} else {
		for _, m := range members.Members {
			if m.UUID != manager.MemberInfo.UUID { // исключение мастера из списка воркеров
				workers.Members = append(workers.Members, m)
			}
		}
	}

	// Проверяем статусы заданий и формируем список актуальных заданий, которые нужно выполнить
	taskInHistory := false // отправить задание в историю
	tasksInQueue := taskpkg.Tasks{}
	for _, t := range tasks.Tasks {
		mbr := manager.Member{UUID: t.Status.WorkerUUID}

		if t.Status.Status == taskpkg.Removing {
			taskInHistory = true
		} else if t.Status.Status == taskpkg.Schedule {
			// Задание распределено, поэтому проверяем статус выполнения

			// Получаем статус задания у самого воркера и решаем что дальше делания с заданием
			if err := getTaskForWorker(db.InstanceETCD, mbr, t); err != nil {
				log.Error("tasksScheduler #6: ", err)
			}

			// Проверяем, что воркер актуальный. Иначе все незавершенные задания перераспределяем между другими воркерами
			okWorker := false
			for _, m := range members.Members {
				if m.UUID == t.Status.WorkerUUID {
					okWorker = true
					break
				}
			}

			if t.Status.Status == taskpkg.Completed {
				taskInHistory = true
			} else if t.Status.Status == taskpkg.Failed {
				if t.NumberOfRetriesOnError > t.Status.RetryCount {
					taskInHistory = false

					if okWorker {
						continue // Не трогаем задание
					}
				} else {
					taskInHistory = true
				}
			} else if t.Status.Status == taskpkg.Pending || t.Status.Status == taskpkg.Running {
				if okWorker {
					// todo: добавить дедлайн на выполнение задания, и если превышение, то отменяем задание.
					continue // Не трогаем задание
				}
			}
		} else if t.Status.Status == taskpkg.Pending {
			taskInHistory = false
		}

		if taskInHistory {
			if err := setTaskInHistory(*t); err != nil {
				log.Error("tasksScheduler #7: ", err)
			}

			// Удаляем задание
			if err := delTaskByProjectETCD(db.InstanceETCD, &mbr, &sources.Project{ID: t.ProjectID}, t); err != nil {
				log.Error("tasksScheduler #8: ", err)
			}

			if t.Status.Status == taskpkg.Removing && t.Status.WorkerUUID != "" {
				// todo: Ключ для прерывания выполнения заданий на воркере. Чтоб успеть отменить еще незапущенные задания
			}
		} else {
			tasksInQueue.Tasks = append(tasksInQueue.Tasks, t)
		}
	}

	if len(tasksInQueue.Tasks) == 0 {
		log.Debug("tasksScheduler #9: no tasks to run")
		return true, nil
	}

	if standalone { // Вариант, когда в кластере только один член, и это мастер
		for _, t := range tasksInQueue.Tasks {
			t.Status.WorkerUUID = manager.MemberInfo.UUID // отмечаем uuid воркера, на который распределено задание
			t.Status.Status = taskpkg.Pending             // для воркера отмечаем статус в ожидании

			ok, err := setTaskForWorker(db.InstanceETCD, manager.MemberInfo, t)
			if err != nil {
				log.Error("tasksScheduler #10: ", err)
			}

			if ok {
				t.Status.Status = taskpkg.Schedule // для общего списка, что задание было распределено наа воркер
			} else {
				t.Status.WorkerUUID = ""
			}
		}
	} else { // Когда есть другие воркеры
		// Распределяем задачи по воркерам
		for _, t := range tasksInQueue.Tasks {
			m, err := GetMemberWithMinTasks(workers)
			if err != nil {
				log.Error("tasksScheduler #11: ", err)
			}

			t.Status.WorkerUUID = m.UUID      // отмечаем uuid воркера, на который распределено задание
			t.Status.Status = taskpkg.Pending // для воркера отмечаем статус в ожидании

			ok, err := setTaskForWorker(db.InstanceETCD, *m, t)
			if err != nil {
				log.Error("tasksScheduler #12: ", err)
			}

			if ok {
				t.Status.Status = taskpkg.Schedule // для общего списка, что задание было распределено наа воркер
			} else {
				t.Status.WorkerUUID = ""
			}
		}
	}

	// Обновление списка заданий по всем проектам
	if err = updateAllTasks(db.InstanceETCD, &tasksInQueue); err != nil {
		log.Error("tasksScheduler #13: ", err)
	}

	return true, nil
}

// setTaskInHistory перенос задания в список истории
func setTaskInHistory(task taskpkg.Task) error {
	log.Debug("setTaskWorker #0: task: ", task.ID)

	// todo: добавить отправку в task_history
	return nil
}

// tasksSchedulerWorker планирование заданий воркером
func tasksSchedulerWorker() (taskpkg.Tasks, error) {
	var tasksInQueue taskpkg.Tasks

	// Получаем список своих задач
	var tasks taskpkg.Tasks
	if err := getTasksForWorker(db.InstanceETCD, manager.MemberInfo, &tasks); err != nil {
		log.Error("tasksSchedulerWorker #0: ", err)
		return tasksInQueue, err
	}

	if len(tasks.Tasks) == 0 {
		log.Debug("tasksSchedulerWorker #1: not found tasks")
		return tasksInQueue, nil
	}

	for _, t := range tasks.Tasks {
		if t.Status.Status == taskpkg.Failed {
			if t.NumberOfRetriesOnError > t.Status.RetryCount {
				tasksInQueue.Tasks = append(tasksInQueue.Tasks, t)
			}
		} else if t.Status.Status == taskpkg.Running {
			// Проверить статус выполнения
			// todo: СДЕЛАТЬ ПРОВЕРКУ ДЕДЛАЙНА И ВСЕ НА ЭТОМ

			log.Info("tasksSchedulerWorker #2: task in Running state (project_id=", t.ProjectID, " job_id=", t.JobID, " task_id=", t.ID, ")")
		} else if t.Status.Status == taskpkg.Pending {
			tasksInQueue.Tasks = append(tasksInQueue.Tasks, t)
		}
	}

	return tasksInQueue, nil
}

// GetMemberWithMinTasks получение члена кластера с минимальным количеством заданий
func GetMemberWithMinTasks(members manager.Members) (*manager.Member, error) {
	memberCountTask := 0
	var member manager.Member
	for _, m := range members.Members {
		var t taskpkg.Tasks
		if err := getTasksForWorker(db.InstanceETCD, m, &t); err != nil {
			log.Error("GetMemberWithMinTasks #0: ", err)
			continue
		}

		if len(t.Tasks) == 0 || len(t.Tasks) < memberCountTask {
			memberCountTask = len(t.Tasks)
			member = m
		}
	}

	return &member, nil
}
