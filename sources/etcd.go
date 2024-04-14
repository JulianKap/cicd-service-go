package sources

import (
	"cicd-service-go/constants"
	"cicd-service-go/db/etcd"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// getProjectsETCD получить список всех проектов из etcd
func getProjectsETCD(cli *clientv3.Client, projects *Projects) error {
	resp, err := etcd.GetKey(cli, Keys.Projects)
	if err != nil {
		log.Error("getProjectsETCD #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getProjectsETCD #1: key ", Keys.Projects, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getProjectsETCD #2: key ", Keys.Projects, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, projects); err != nil {
		log.Error("getProjectsETCD #3: ", err)
		return err
	}

	return nil
}

// getProjectETCD получить указанный проект из etcd по ID
func getProjectETCD(cli *clientv3.Client, project *Project) (bool, error) {
	key := Keys.Project + strconv.Itoa(project.ID)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getProjectETCD #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getProjectETCD #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getProjectETCD #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &project); err != nil {
		log.Error("getProjectETCD #3: ", err)
		return false, err
	}

	return true, nil
}

// createProjectETCD создание нового проекта
func (p *Project) createProjectETCD(cli *clientv3.Client) error {
	latestId, err := getLatestIdETCD(cli, Keys.LatestID)
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
	projects.Projects = append(projects.Projects, *p)

	projectsJSON, err := json.Marshal(projects)
	if err != nil {
		log.Error("createProjectETCD #2: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Projects, string(projectsJSON)); err != nil {
		log.Error("createProjectETCD #3: ", err)
		return err
	}

	// Добавление проекта в отдельный ключ
	projectJSON, err := json.Marshal(p)
	if err != nil {
		log.Error("createProjectETCD #4: ", err)
		return err
	}

	if err = etcd.SetKey(cli, Keys.Project+strconv.Itoa(p.ID), string(projectJSON)); err != nil {
		log.Error("createProjectETCD #5: ", err)
		return err
	}

	// Добавление последнего ID
	if err = etcd.SetKey(cli, Keys.LatestID, strconv.Itoa(p.ID)); err != nil {
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
		log.Info("deleteProjectETCD #2: this project has ", len(jobs.Jobs), " tasks. need to remove them first")
		return false, "Error delete project. Has " + strconv.Itoa(len(jobs.Jobs)) + " tasks. Need to remove them first", nil
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
	if err = etcd.SetKey(cli, Keys.Projects, string(valueJSON)); err != nil {
		log.Error("deleteProjectETCD #4: ", err)
		return state, "Error update projects key", err
	}

	// Удаляем проект
	key := Keys.Project + strconv.Itoa(p.ID)
	if err = etcd.DelKey(cli, key); err != nil {
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
func (p *Project) getJobsETCD(cli *clientv3.Client, jobs *Jobs) error {
	key := Keys.Project + strconv.Itoa(p.ID) + constants.JOBS
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getJobs #0: ", err)
		return err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getJobs #1: key ", key, " not found")
		return nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getJobs #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &jobs); err != nil {
		log.Error("getJobs #3: ", err)
		return err
	}

	return nil
}

// getJobETCD получить задачу
func (p *Project) getJobETCD(cli *clientv3.Client, job *Job) (bool, error) {
	key := Keys.Project + strconv.Itoa(p.ID) + constants.JOBS + strconv.Itoa(job.ID)
	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getJob #0: ", err)
		return false, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getJob #1: key ", key, " not found")
		return false, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getJob #2: key ", key, " get more than one key")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &job); err != nil {
		log.Error("getJob #3: ", err)
		return false, err
	}

	return true, nil
}

// createJobETCD добавление задачи
func (p *Project) createJobETCD(cli *clientv3.Client, j *Job) error {
	keyLatestId := Keys.Project + strconv.Itoa(p.ID) + constants.JOB_LATEST_ID
	latestId, err := getLatestIdETCD(cli, keyLatestId)
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
	jobs.Jobs = append(jobs.Jobs, *j)

	jobsJSON, err := json.Marshal(jobs)
	if err != nil {
		log.Error("createJob #2: ", err)
		return err
	}

	key := Keys.Project + strconv.Itoa(p.ID) + constants.JOBS
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

	key = Keys.Project + strconv.Itoa(p.ID) + constants.JOBS + strconv.Itoa(j.ID)
	if err = etcd.SetKey(cli, key, string(jobJSON)); err != nil {
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
	key := Keys.Project + strconv.Itoa(p.ID) + constants.JOBS
	if err = etcd.SetKey(cli, key, string(valueJSON)); err != nil {
		log.Error("deleteJobETCD #2: ", err)
		return state, err
	}

	// Удаляем задачу
	key = Keys.Project + strconv.Itoa(p.ID) + constants.JOBS + strconv.Itoa(j.ID)
	if err = etcd.DelKey(cli, key); err != nil {
		log.Error("deleteJobETCD #3: ", err)
		return state, err
	}

	return state, nil
}

// getLatestIdETCD получить последний id проекта (или задачи)
func getLatestIdETCD(cli *clientv3.Client, key string) (int, error) {
	id := -1

	resp, err := etcd.GetKey(cli, key)
	if err != nil {
		log.Error("getLatestProjectID #0: ", err)
		return id, err
	}

	// Проверка наличия ключа
	if len(resp.Kvs) == 0 { // Ключ не найден
		log.Info("getLatestProjectID #1: key ", key, " not found")
		return id, nil
	} else if len(resp.Kvs) > 1 { // Больше одного ключа
		log.Warning("getLatestProjectID #2: key ", key, " get more than one key")
	}

	value := string(resp.Kvs[0].Value)
	id, err = strconv.Atoi(value)
	if err != nil {
		log.Error("getLatestProjectID #3: ", err)
		return id, err
	}

	return id, nil
}
