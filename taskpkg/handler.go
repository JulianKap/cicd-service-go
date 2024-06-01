package taskpkg

import (
	"cicd-service-go/init/db"
	"cicd-service-go/sources"
	"cicd-service-go/utility"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// HandleTaskCreate создание задания
func HandleTaskCreate(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var task Task
	if err := ctx.Bind(&task); err != nil {
		log.Error("HandleTaskCreate #0: ", err)
		return ctx.JSON(http.StatusBadRequest, TaskResponse{
			Message: "Error convert struct task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	//todo: оформить заполнение данных и валидацию в отдельной функции
	tm := time.Now()
	task.CreateAt = &tm

	// Проверка существования проекта
	p := sources.Project{ID: task.ProjectID}
	codeValPrj, respValPrj := sources.ValidateProject(&p)
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	// Проверка существования задачи
	j := sources.Job{ID: task.JobID}
	codeValJob, respValJob := sources.ValidateJob(&p, &j)
	if codeValJob != http.StatusOK {
		return ctx.JSON(codeValJob, respValJob)
	}

	if err := task.setTaskByProjectETCD(db.InstanceETCD, &p); err != nil {
		log.Error("HandleTaskCreate #1: ", err)
		return ctx.JSON(http.StatusBadRequest, TasksResponse{
			Message: "Error create task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, TaskResponse{
		Task:    &task,
		Message: "Task create",
	})
}

// HandleTasksGetList получить список всех заданий
func HandleTasksGetList(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	var tasks Tasks
	if err := tasks.getTasksByProjectETCD(db.InstanceETCD, &project, false); err != nil {
		log.Error("HandleTasksGetList #0: ", err)
		return ctx.JSON(http.StatusInternalServerError, TasksResponse{
			Message: "Error get tasks list",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, TasksResponse{
		Tasks: tasks.Tasks,
	})
}

// HandleTaskGetByID получить конкретное задание по id
func HandleTaskGetByID(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	taskID, err := strconv.Atoi(ctx.Param("id_task"))
	if err != nil {
		log.Error("HandleTaskGetByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, TaskResponse{
			Message: "Error convert id task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	task := Task{
		ID:        taskID,
		ProjectID: project.ID,
	}

	state, err := task.getTaskByProjectETCD(db.InstanceETCD)
	if err != nil {
		log.Error("HandleTaskGetByID #1: ", err)
		return ctx.JSON(http.StatusInternalServerError, TaskResponse{
			Message: "Error get task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, TaskResponse{
			Task: &task,
		})
	} else {
		log.Info("HandleTaskGetByID #2: not found task")
		return ctx.JSON(http.StatusBadRequest, TaskResponse{
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
	codeValPrj, respValPrj := sources.ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	taskID, err := strconv.Atoi(ctx.Param("id_task"))
	if err != nil {
		log.Error("HandleTaskDeleteByID #0: ", err)
		return ctx.JSON(http.StatusBadRequest, TaskResponse{
			Message: "Error convert id task",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	task := Task{
		ID:        taskID,
		ProjectID: project.ID,
	}

	state, err := task.markerDelTaskByProjectETCD(db.InstanceETCD, &project)
	if err != nil {
		log.Error("HandleTaskDeleteByID #1: ", err)
		return ctx.JSON(http.StatusBadRequest, TaskResponse{
			Error: utility.StringPtr(err.Error()),
		})
	}

	if state {
		return ctx.JSON(http.StatusOK, TaskResponse{Message: "Task delete"})
	} else {
		log.Info("HandleTaskDeleteByID #2: not found task")
		return ctx.JSON(http.StatusBadRequest, TaskResponse{Message: "Task not found"})
	}
}

// HandleTasksGetHistoryList получить список историчных заданий
func HandleTasksGetHistoryList(ctx echo.Context) (err error) {
	codeValidation, respValidation := sources.ValidatePermission()
	if codeValidation != http.StatusOK {
		return ctx.JSON(codeValidation, respValidation)
	}

	var project sources.Project
	codeValPrj, respValPrj := sources.ValidateProjectById(ctx, &project, "id_project")
	if codeValPrj != http.StatusOK {
		return ctx.JSON(codeValPrj, respValPrj)
	}

	var tasks Tasks
	if err := tasks.getTasksByProjectETCD(db.InstanceETCD, &project, true); err != nil {
		log.Error("HandleTasksGetHistoryList #0: ", err)
		return ctx.JSON(http.StatusInternalServerError, TasksResponse{
			Message: "Error get tasks list",
			Error:   utility.StringPtr(err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, TasksResponse{
		Tasks: tasks.Tasks,
	})
}
