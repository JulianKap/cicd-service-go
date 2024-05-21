package sources

import (
	"cicd-service-go/constants"
	"cicd-service-go/db/etcd"
	"cicd-service-go/manager"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GetKeyProjects получить ключ всех проектов
func GetKeyProjects() string {
	return manager.Conf.Cluster.Namespace + constants.PROJECTS_ALL
}

// GetKeyLatestIdProject получить ключ идентификатора последнего проекта
func GetKeyLatestIdProject() string {
	return manager.Conf.Cluster.Namespace + constants.PROJECT_LATEST_ID
}

// GetKeyProject получить ключ конкретного проекта
func GetKeyProject(p *Project) string {
	return manager.Conf.Cluster.Namespace + constants.PROJECTS + "/" + strconv.Itoa(p.ID)
}

// GetKeyJobsByProject получить ключ всех задач конкретного проекта
func GetKeyJobsByProject(p *Project) string {
	return GetKeyProject(p) + constants.JOBS
}

// GetKeyJobLatestIdProject получить ключ идентификатора последней задачи конкретного проекта
func GetKeyJobLatestIdProject(p *Project) string {
	return GetKeyProject(p) + constants.JOB_LATEST_ID
}

// GetKeyJobByProject получить ключ конкретной задачи конкретного проекта
func GetKeyJobByProject(p *Project, j *Job) string {
	return GetKeyJobsByProject(p) + "/" + strconv.Itoa(j.ID)
}

// getProjectsETCD получить список всех проектов из etcd
func getProjectsETCD(cli *clientv3.Client, p *Projects) error {
	key := GetKeyProjects()
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getProjectsETCD #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getProjectsETCD #1: key ", key, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getProjectsETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, p); err != nil {
		log.Error("getProjectsETCD #3: ", err)
		return err
	}

	return nil
}

// getProjectETCD получить указанный проект из etcd по ID
func getProjectETCD(cli *clientv3.Client, p *Project) (bool, error) {
	key := GetKeyProject(p)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getProjectETCD #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getProjectETCD #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getProjectETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &p); err != nil {
		log.Error("getProjectETCD #3: ", err)
		return false, err
	}

	return true, nil
}

// createProjectETCD создание нового проекта
func (p *Project) createProjectETCD(cli *clientv3.Client) error {
	latestId, err := etcd.GetKeyInt(cli, GetKeyLatestIdProject())
	if err != nil {
		log.Error("createProjectETCD #0: ", err)
		return err
	}

	if latestId == -1 {
		latestId = 1
	}
	p.ID = latestId + 1

	// Добавление проекта в список всех проектов
	var projects Projects
	if err := getProjectsETCD(cli, &projects); err != nil {
		log.Error("createProjectETCD #1: ", err)
		return err
	}
	projects.Projects = append(projects.Projects, p)

	projectsJSON, err := json.Marshal(projects)
	if err != nil {
		log.Error("createProjectETCD #2: ", err)
		return err
	}

	if err = etcd.SetKey(cli, GetKeyProjects(), string(projectsJSON)); err != nil {
		log.Error("createProjectETCD #3: ", err)
		return err
	}

	// Добавление проекта в отдельный ключ
	projectJSON, err := json.Marshal(p)
	if err != nil {
		log.Error("createProjectETCD #4: ", err)
		return err
	}

	if err = etcd.SetKey(cli, GetKeyProject(p), string(projectJSON)); err != nil {
		log.Error("createProjectETCD #5: ", err)
		return err
	}

	// Добавление последнего ID
	if err = etcd.SetKey(cli, GetKeyLatestIdProject(), strconv.Itoa(p.ID)); err != nil {
		log.Error("createProjectETCD #6: ", err)
		return err
	}

	return nil
}

// deleteProjectETCD удалить проект
func (p *Project) deleteProjectETCD(cli *clientv3.Client) (bool, string, error) {
	var projects Projects
	if err := getProjectsETCD(cli, &projects); err != nil {
		log.Error("deleteProjectETCD #0: ", err)
		return false, "Error get project", err
	}

	var jobs Jobs
	if err := p.getJobsETCD(cli, &jobs); err != nil {
		log.Error("deleteProjectETCD #1: ", err)
		return false, "Error get jobs list for this project", err
	}

	if len(jobs.Jobs) > 0 {
		log.Info("deleteProjectETCD #2: this project has ", len(jobs.Jobs), " jobs. Need to remove them first")
		return false, "Error delete project. Has " + strconv.Itoa(len(jobs.Jobs)) + " jobs. Need to remove them first", nil
	}

	state := false
	var newProjects Projects
	for _, project := range projects.Projects {
		if p.ID != project.ID {
			newProjects.Projects = append(newProjects.Projects, project)
		} else {
			state = true
		}
	}

	valueJSON, err := json.Marshal(newProjects)
	if err != nil {
		log.Error("deleteProjectETCD #3: ", err)
		return state, "Error encoding projects", err
	}

	// Обновляем список всех проектов
	if err = etcd.SetKey(cli, GetKeyProjects(), string(valueJSON)); err != nil {
		log.Error("deleteProjectETCD #4: ", err)
		return state, "Error update projects key", err
	}

	// Удаляем проект
	if err = etcd.DelKey(cli, GetKeyProject(p)); err != nil {
		log.Error("deleteProjectETCD #5: ", err)
		return state, "Error delete project key", err
	}

	if state {
		return state, "Project delete", nil
	} else {
		return state, "Project not found", nil
	}
}

// getJobsETCD получить список всех задач проекта
func (p *Project) getJobsETCD(cli *clientv3.Client, j *Jobs) error {
	key := GetKeyJobsByProject(p)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getJobs #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getJobs #1: key ", key, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getJobs #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &j); err != nil {
		log.Error("getJobs #3: ", err)
		return err
	}

	return nil
}

// GetJobETCD получить задачу
func (p *Project) GetJobETCD(cli *clientv3.Client, j *Job) (bool, error) {
	key := GetKeyJobByProject(p, j)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getJob #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Debug("getJob #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getJob #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &j); err != nil {
		log.Error("getJob #3: ", err)
		return false, err
	}

	return true, nil
}

// createJobETCD добавление задачи
func (p *Project) createJobETCD(cli *clientv3.Client, j *Job) error {
	keyLatestId := GetKeyJobLatestIdProject(p)
	latestId, err := etcd.GetKeyInt(cli, keyLatestId)
	if err != nil {
		log.Error("createJob #0: ", err)
		return err
	}

	if latestId == -1 {
		latestId = 1
	}
	j.ID = latestId + 1

	// Добавление задачи в список всех задач проекта
	var jobs Jobs
	if err := p.getJobsETCD(cli, &jobs); err != nil {
		log.Error("createJob #1: ", err)
		return err
	}
	jobs.Jobs = append(jobs.Jobs, j)

	jobsJSON, err := json.Marshal(jobs)
	if err != nil {
		log.Error("createJob #2: ", err)
		return err
	}

	key := GetKeyJobsByProject(p)
	if err = etcd.SetKey(cli, key, string(jobsJSON)); err != nil {
		log.Error("createJob #3: ", err)
		return err
	}

	// Добавление проекта в отдельный ключ
	jobJSON, err := json.Marshal(j)
	if err != nil {
		log.Error("createJob #4: ", err)
		return err
	}

	if err = etcd.SetKey(cli, GetKeyJobByProject(p, j), string(jobJSON)); err != nil {
		log.Error("createJob #5: ", err)
		return err
	}

	// Добавление последнего ID
	if err = etcd.SetKey(cli, keyLatestId, strconv.Itoa(j.ID)); err != nil {
		log.Error("createJob #6: ", err)
		return err
	}

	return nil
}

// deleteJobETCD удаление задачи
func (p *Project) deleteJobETCD(cli *clientv3.Client, j *Job) (bool, error) {
	var jobs Jobs
	if err := p.getJobsETCD(cli, &jobs); err != nil {
		log.Error("deleteJobETCD #0: ", err)
		return false, err
	}

	state := false
	var newJobs Jobs
	for _, job := range jobs.Jobs {
		if j.ID != job.ID {
			newJobs.Jobs = append(newJobs.Jobs, job)
		} else {
			state = true
		}
	}

	valueJSON, err := json.Marshal(newJobs)
	if err != nil {
		log.Error("deleteJobETCD #1: ", err)
		return state, err
	}

	// Обновляем список всех задач
	if err = etcd.SetKey(cli, GetKeyJobsByProject(p), string(valueJSON)); err != nil {
		log.Error("deleteJobETCD #2: ", err)
		return state, err
	}

	// Удаляем задачу
	if err = etcd.DelKey(cli, GetKeyJobByProject(p, j)); err != nil {
		log.Error("deleteJobETCD #3: ", err)
		return state, err
	}

	return state, nil
}
