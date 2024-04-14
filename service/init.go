package service

import (
	"cicd-service-go/schedule"
	"cicd-service-go/sources"
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
	schedule.InitHandler()

	// Проекты
	project := e.Group("/project")
	project.PUT("/create", sources.HandleProjectCreate)
	project.GET("/all", sources.HandleProjectsGetList)
	project.GET("/:id", sources.HandleProjectGetByID)
	project.DELETE("/:id", sources.HandleProjectDeleteByID)

	// Задачи
	jobs := project.Group("/jobs")
	jobs.PUT("/create", sources.HandleJobCreate)
	jobs.GET("/:id/all", sources.HandleJobsGetList)
	jobs.GET("/:id_project/:id_job", sources.HandleJobGetByID)
	jobs.DELETE("/:id_project/:id_job", sources.HandleJobDeleteByID)

	// todo: сделать роуты для обновления проектов, задач
	// В частности обновление названий, токенов, кредов

	// Таски
	tasks := project.Group("/tasks")
	tasks.PUT("/create", schedule.HandleTaskCreate)
	tasks.GET("/:id/all", schedule.HandleTasksGetList)
	tasks.GET("/:id_project/:id_task", schedule.HandleTaskGetByID)
	tasks.DELETE("/:id_project/:id_task", schedule.HandleTaskDeleteByID)

	e.GET("/hc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
}
