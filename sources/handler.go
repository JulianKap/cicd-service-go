package sources

import (
	"cicd-service-go/init/db"
	"cicd-service-go/init/secrets"
	"cicd-service-go/manager"
	"cicd-service-go/utility"
	"cicd-service-go/vault"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func HandleProjectCreate(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	var p Project
	if err := ctx.Bind(&p); err != nil {
		log.Error("HandleProjectCreate #0: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error convert struct project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if p.ProjectName == "" {
		log.Info("HandleProjectCreate #1: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Empty project name",
		})
	}

	var projects Projects
	if err := getProjectsETCD(db.InstanceETCD, &projects); err != nil {
		log.Error("HandleProjectCreate #2: ", err)
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "Error get projects list for validate",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	for _, project := range projects.Projects {
		if p.ProjectName == project.ProjectName {
			log.Info("HandleProjectCreate #3: project name already exists")
			return ctx.JSON(http.StatusBadRequest, Response{
				Message: "Project name already exists",
			})
		}
	}

	var t vault.Token
	if err := p.createTokenProjectVault(secrets.InstanceVault, &t); err != nil {
		log.Error("HandleProjectCreate #4: ", err)
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "Error generate access token",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if err = p.createProjectETCD(db.InstanceETCD); err != nil {
		log.Error("HandleProjectCreate #5: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error create project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	p.Token = t.Token
	return ctx.JSON(http.StatusOK, ProjectResponse{
		Message: "Project created",
		Project: &p,
	})
}

// HandleProjectsGetList получить список всех проектов
func HandleProjectsGetList(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	var projects Projects
	if err := getProjectsETCD(db.InstanceETCD, &projects); err != nil {
		log.Error("HandleProjectsGetList #0: ", err)
		return ctx.JSON(http.StatusInternalServerError, ProjectResponse{
			Message: "Error get projects list",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	for _, p := range projects.Projects {
		var t vault.Token
		t.Path = GetProjectPathToken(p)
		if err := vault.GetToken(secrets.InstanceVault, &t); err != nil {
			log.Error("HandleProjectsGetList #1: ", err)
		}
		p.Token = t.Token
	}

	return ctx.JSON(http.StatusOK, ProjectsResponse{
		Projects: projects.Projects,
	})
}

// HandleProjectGetByID вывести статус по id
func HandleProjectGetByID(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	projectID, err := strconv.Atoi(ctx.Param("id_project"))
	if err != nil {
		log.Error("HandleProjectGetByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	p := Project{
		ID: projectID,
	}
	state, err := getProjectETCD(db.InstanceETCD, &p)
	if err != nil {
		log.Error("HandleProjectGetByID #1: ", err)
		return ctx.JSON(http.StatusInternalServerError, ProjectResponse{
			Message: "Error get project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	var t vault.Token
	t.Path = GetProjectPathToken(&p)
	if err := vault.GetToken(secrets.InstanceVault, &t); err != nil {
		log.Error("HandleProjectGetByID #2: ", err)
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "Not found token for project " + p.ProjectName,
		})
	}

	p.Token = t.Token
	if state {
		return ctx.JSON(http.StatusOK, ProjectResponse{
			Project: &p,
		})
	} else {
		log.Info("HandleProjectGetByID #3: not found project")
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Not found project",
		})
	}
}

// HandleProjectDeleteByID удаление проекта по id
func HandleProjectDeleteByID(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	projectID, err := strconv.Atoi(ctx.Param("id_project"))
	if err != nil {
		log.Error("HandleProjectDeleteByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Error: utility.StringPtr(err.Error()),
		})
	}

	project := Project{
		ID: projectID,
	}
	state, message, err := project.deleteProjectETCD(db.InstanceETCD)
	if err != nil {
		log.Error("HandleProjectDeleteByID #1: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Error: utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, Response{Message: message})
	} else {
		log.Info("HandleProjectDeleteByID #2: not found project")
		return ctx.JSON(http.StatusBadRequest, Response{Message: message})
	}
}

// HandleJobCreate создать задачу
func HandleJobCreate(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	var job Job
	if err := ctx.Bind(&job); err != nil {
		log.Error("HandleJobCreate #0: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error convert struct job",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if job.JobName == "" {
		log.Info("HandleJobCreate #1: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Empty job name",
		})
	}

	project := Project{ID: job.IdProject}
	codeValPrj, respValPrj := ValidateProject(&project)
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	var jobs Jobs
	if err := project.getJobsETCD(db.InstanceETCD, &jobs); err != nil {
		log.Error("HandleJobCreate #4: ", err)
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "Error get jobs list for validate",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	for _, j := range jobs.Jobs {
		if j.JobName == job.JobName {
			log.Info("HandleJobCreate #5: job name already exists")
			return ctx.JSON(http.StatusBadRequest, Response{
				Message: "Job name already exists",
			})
		}
	}

	if err = project.createJobETCD(db.InstanceETCD, &job); err != nil {
		log.Error("HandleJobCreate #6: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error create job",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, JobResponse{
		Message: "Job created",
		Job:     &job,
	})
}

// HandleJobsGetList получить список задач проекта
func HandleJobsGetList(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	var project Project
	codeValPrj, respValPrj := ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	var jobs Jobs
	if err := project.getJobsETCD(db.InstanceETCD, &jobs); err != nil {
		log.Error("HandleJobsGetList #0: ", err)
		return ctx.JSON(http.StatusInternalServerError, ProjectResponse{
			Message: "Error get jobs list",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, JobsResponse{
		Jobs: jobs.Jobs,
	})
}

// HandleJobGetByID получить задачу проекта
func HandleJobGetByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project Project
	codeValPrj, respValPrj := ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	jobID, err := strconv.Atoi(ctx.Param("id_job"))
	if err != nil {
		log.Error("HandleJobGetByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id job",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	job := Job{
		ID: jobID,
	}

	state, err := project.GetJobETCD(db.InstanceETCD, &job)
	if err != nil {
		log.Error("HandleJobGetByID #1: ", err)
		return ctx.JSON(http.StatusInternalServerError, ProjectResponse{
			Message: "Error get job",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, JobResponse{
			Job: &job,
		})
	} else {
		log.Info("HandleJobGetByID #2: not found job")
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Not found job",
		})
	}
}

// HandleJobDeleteByID удалить задачу по id
func HandleJobDeleteByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project Project
	codeValPrj, respValPrj := ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	jobID, err := strconv.Atoi(ctx.Param("id_job"))
	if err != nil {
		log.Error("HandleJobDeleteByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id job",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	job := Job{
		ID: jobID,
	}

	state, err := project.deleteJobETCD(db.InstanceETCD, &job)
	if err != nil {
		log.Error("HandleJobDeleteByID #1: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Error: utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, Response{Message: "Job delete"})
	} else {
		log.Info("HandleJobDeleteByID #2: not found job")
		return ctx.JSON(http.StatusBadRequest, Response{Message: "Job not found"})
	}
}

// ValidatePermission валидация доступных прав для данного запущенного экземпляра сервиса
func ValidatePermission() (int, Response) {
	// Если мы не мастер, то отклоняем данный запрос
	if manager.MemberInfo.Role != manager.MasterRole {
		log.Info("validatePermission #0: not master")
		return http.StatusMethodNotAllowed, Response{Message: "I am SLAVE! Slave not support management projects and jobs"}
	}

	// Если в режиме ro, то отклоняем данный запрос
	if manager.MemberInfo.ReadOnly {
		log.Info("validatePermission #1: in read only state")
		return http.StatusMethodNotAllowed, Response{Message: "My state is read only"}
	}

	return http.StatusOK, Response{}
}

// ValidateProject проверка существования проекта и получение его
func ValidateProject(p *Project) (int, ProjectResponse) {
	state, err := getProjectETCD(db.InstanceETCD, p)
	if err != nil {
		log.Error("ValidateProject #0: ", err)
		return http.StatusInternalServerError, ProjectResponse{
			Message: "Error get project",
			Error:   utility.StringPtr(err.Error()),
		}
	}

	if !state {
		log.Info("ValidateProject #1: ", err)
		return http.StatusBadRequest, ProjectResponse{
			Message: "Not found project",
		}
	}

	return http.StatusOK, ProjectResponse{}
}

// ValidateProjectById проверка существования проекта и получение его по id из запроса
func ValidateProjectById(ctx echo.Context, p *Project, id string) (int, ProjectResponse) {
	projectID, err := strconv.Atoi(ctx.Param(id))
	if err != nil {
		log.Error("ValidateProjectById #0: ", err)
		return http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id project",
			Error:   utility.StringPtr(err.Error()),
		}
	}

	p.ID = projectID
	return ValidateProject(p)
}

// ValidateJob проверка существования задачи и получение ее
func ValidateJob(p *Project, j *Job) (int, JobResponse) {
	state, err := p.GetJobETCD(db.InstanceETCD, j)
	if err != nil {
		log.Error("ValidateJob #0: ", err)
		return http.StatusInternalServerError, JobResponse{
			Message: "Error get job",
			Error:   utility.StringPtr(err.Error()),
		}
	}

	if !state {
		log.Info("ValidateJob #1: ", err)
		return http.StatusBadRequest, JobResponse{
			Message: "Not found job",
		}
	}

	return http.StatusOK, JobResponse{}
}
