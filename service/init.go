package service

import (
	"cicd-service-go/manager"
	"cicd-service-go/sources"
	"cicd-service-go/taskpkg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func startFramework() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	initRoutes(e)
	return e
}

func initRoutes(e *echo.Echo) {
	sources.InitHandler()
	taskpkg.InitHandler()

	// Projects
	project := e.Group("/project")
	project.PUT("/create", sources.HandleProjectCreate)
	project.GET("/all", sources.HandleProjectsGetList)
	project.GET("/:id", sources.HandleProjectGetByID)
	project.DELETE("/:id", sources.HandleProjectDeleteByID)

	// Jobs
	jobs := project.Group("/jobs")
	jobs.PUT("/create", sources.HandleJobCreate)
	jobs.GET("/:id/all", sources.HandleJobsGetList)
	jobs.GET("/:id_project/:id_job", sources.HandleJobGetByID)
	jobs.DELETE("/:id_project/:id_job", sources.HandleJobDeleteByID)

	// todo: сделать роуты для обновления проектов, задач
	// В частности обновление названий, токенов, кредов

	// Tasks
	tasks := project.Group("/tasks")
	tasks.PUT("/create", taskpkg.HandleTaskCreate)
	tasks.GET("/:id/all", taskpkg.HandleTasksGetList)
	tasks.GET("/:id_project/:id_task", taskpkg.HandleTaskGetByID)
	tasks.DELETE("/:id_project/:id_task", taskpkg.HandleTaskDeleteByID)

	// Проверка на мастера
	e.GET("/master", func(c echo.Context) error {
		if manager.MemberInfo.Role == manager.MasterRole {
			return c.JSON(http.StatusOK, manager.MemberInfo)
		}
		return c.JSON(http.StatusBadRequest, manager.MemberInfo)
	})

	// Healthcheck
	e.GET("/hc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
}
