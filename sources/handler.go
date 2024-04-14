package sources

import (
	"cicd-service-go/constants"
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"cicd-service-go/utility"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Keys KeysDCS
)

func InitHandler() {
	namespace := viper.GetString("cluster.namespace_dcs")
	Keys = KeysDCS{
		Projects: namespace + constants.PROJECTS_ALL,
		LatestID: namespace + constants.PROJECT_LATEST_ID,
		Project:  namespace + constants.PROJECTS + "/",
	}
}

func HandleProjectCreate(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	var project Project
	if err := ctx.Bind(&project); err != nil {
		log.Error("HandleProjectCreate #0: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error convert struct project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if project.ProjectName == "" {
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

	for _, p := range projects.Projects {
		if project.ProjectName == p.ProjectName {
			log.Info("HandleProjectCreate #3: project name already exists")
			return ctx.JSON(http.StatusBadRequest, Response{
				Message: "Project name already exists",
			})
		}
	}

	token, err := utility.GenerateToken(16)
	if err != nil {
		log.Error("HandleProjectCreate #4: ", err)
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "Error generate access token",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	project.APIKey = token
	if err = project.createProjectETCD(db.InstanceETCD); err != nil {
		log.Error("HandleProjectCreate #5: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{
			Message: "Error create project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, ProjectResponse{
		Message: "Project created",
		Project: &project,
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

	return ctx.JSON(http.StatusOK, ProjectsResponse{
		Projects: &projects,
	})
}

// HandleProjectGetByID вывести статус по id
func HandleProjectGetByID(ctx echo.Context) (err error) {
	codeValPerm, respValPerm := ValidatePermission()
	if codeValPerm != http.StatusOK {
		return ctx.JSON(codeValPerm, respValPerm)
	}

	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error("HandleProjectGetByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	project := Project{
		ID: projectID,
	}
	state, err := getProjectETCD(db.InstanceETCD, &project)
	if err != nil {
		log.Error("HandleProjectGetByID #1: ", err)
		return ctx.JSON(http.StatusInternalServerError, ProjectResponse{
			Message: "Error get project",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, ProjectResponse{
			Project: &project,
		})
	} else {
		log.Info("HandleProjectGetByID #2: not found project")
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

	projectID, err := strconv.Atoi(ctx.Param("id"))
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

	var project Project
	codeValPrj, respValPrj := ValidateProjectById(&project, job.IdProject)
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
	codeValPrj, respValPrj := ValidateProjectByIdContext(ctx, &project, "id")
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
		Jobs: &jobs,
	})
}

// HandleJobGetByID получить задачу проекта
func HandleJobGetByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project Project
	codeValPrj, respValPrj := ValidateProjectByIdContext(ctx, &project, "id_project")
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

	state, err := project.getJobETCD(db.InstanceETCD, &job)
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
	codeValPrj, respValPrj := ValidateProjectByIdContext(ctx, &project, "id_project")
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
	if !manager.MemberInfo.Master {
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

// ValidateProjectByIdContext проверка существования проекта и получение его по id из запроса
func ValidateProjectByIdContext(ctx echo.Context, project *Project, id string) (int, ProjectResponse) {
	projectID, err := strconv.Atoi(ctx.Param(id))
	if err != nil {
		log.Error("validateProjectByIdContext #0: ", err)
		return http.StatusBadRequest, ProjectResponse{
			Message: "Error convert id project",
			Error:   utility.StringPtr(err.Error()),
		}
	}

	return ValidateProjectById(project, projectID)
}

// ValidateProjectById проверка существования проекта и получение его
func ValidateProjectById(project *Project, id int) (int, ProjectResponse) {
	project.ID = id
	state, err := getProjectETCD(db.InstanceETCD, project)
	if err != nil {
		log.Error("validateProjectById #0: ", err)
		return http.StatusInternalServerError, ProjectResponse{
			Message: "Error get project",
			Error:   utility.StringPtr(err.Error()),
		}
	}

	if !state {
		log.Info("validateProjectById #1: ", err)
		return http.StatusBadRequest, ProjectResponse{
			Message: "Not found project",
		}
	}

	return http.StatusOK, ProjectResponse{}
}
