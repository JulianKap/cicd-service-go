package schedule

import (
	"cicd-service-go/constants"
	"cicd-service-go/init/db"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"cicd-service-go/utility"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

var (
	Keys taskpkg.KeysDCS
)

func InitHandler() {
	namespace := viper.GetString("cluster.namespace_dcs")
	Keys = taskpkg.KeysDCS{
		Tasks:        namespace + constants.PROJECTS_TASKS,
		TasksHistory: namespace + constants.PROJECTS_TASKS_HISTORY,
	}
}

// HandleTaskCreate создание таски
func HandleTaskCreate(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var task taskpkg.Task
	if err := ctx.Bind(&task); err != nil {
		log.Error("HandleTaskCreate #0: ", err)
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{
			Message: "Error convert struct task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectById(&project, task.ProjectID)
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	//var tasks taskpkg.Tasks
	//if err := getTasksETCD(db.InstanceETCD, &project, &tasks); err != nil {
	//	log.Error("HandleTaskCreate #1: ", err)
	//	return ctx.JSON(http.StatusInternalServerError, taskpkg.TasksResponse{
	//		Message: "Error get tasks list",
	//		Error:   utility.StringPtr(err.Error()),
	//	})
	//}

	return nil
}

// HandleTasksGetList получить список всех тасок
func HandleTasksGetList(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectByIdContext(ctx, &project, "id")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	var tasks taskpkg.Tasks
	if err := getTasksETCD(db.InstanceETCD, &project, &tasks); err != nil {
		log.Error("HandleTasksGetList #0: ", err)
		return ctx.JSON(http.StatusInternalServerError, taskpkg.TasksResponse{
			Message: "Error get tasks list",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, taskpkg.TasksResponse{
		Tasks: &tasks,
	})
}

// HandleTaskGetByID получить конкретную таску по id
func HandleTaskGetByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectByIdContext(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	taskID, err := strconv.Atoi(ctx.Param("id_task"))
	if err != nil {
		log.Error("HandleTaskGetByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{
			Message: "Error convert id task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	task := taskpkg.Task{
		ID: taskID,
	}

	state, err := getTaskETCD(db.InstanceETCD, &project, &task)
	if err != nil {
		log.Error("HandleTaskGetByID #1: ", err)
		return ctx.JSON(http.StatusInternalServerError, taskpkg.TaskResponse{
			Message: "Error get task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, taskpkg.TaskResponse{
			Task: &task,
		})
	} else {
		log.Info("HandleTaskGetByID #2: not found task")
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{
			Message: "Not found task",
		})
	}
}

// HandleTaskDeleteByID удалить таксу по id
func HandleTaskDeleteByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectByIdContext(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	taskID, err := strconv.Atoi(ctx.Param("id_task"))
	if err != nil {
		log.Error("HandleTaskDeleteByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{
			Message: "Error convert id task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	task := taskpkg.Task{
		ID: taskID,
	}

	state, err := deleteTaskETCD(db.InstanceETCD, &project, &task)
	if err != nil {
		log.Error("HandleTaskDeleteByID #1: ", err)
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{
			Error: utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, taskpkg.TaskResponse{Message: "Task delete"})
	} else {
		log.Info("HandleTaskDeleteByID #2: not found task")
		return ctx.JSON(http.StatusBadRequest, taskpkg.TaskResponse{Message: "Task not found"})
	}
}
