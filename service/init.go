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
	allProjectsAuthMiddleware := allProjectsAuthMiddleware()
	projectAuthMiddleware := projectAuthMiddleware()

	// Projects
	project := e.Group("/project")
	project.PUT("/create", sources.HandleProjectCreate, allProjectsAuthMiddleware)
	project.GET("/all", sources.HandleProjectsGetList, allProjectsAuthMiddleware)
	project.GET("/:id_project", sources.HandleProjectGetByID, allProjectsAuthMiddleware)
	project.DELETE("/:id_project", sources.HandleProjectDeleteByID, allProjectsAuthMiddleware)

	// Jobs
	jobs := project.Group("/jobs")
	jobs.PUT("/:id_project/create", sources.HandleJobCreate, projectAuthMiddleware)
	jobs.GET("/:id_project/all", sources.HandleJobsGetList, projectAuthMiddleware)
	jobs.GET("/:id_project/:id_job", sources.HandleJobGetByID, projectAuthMiddleware)
	jobs.DELETE("/:id_project/:id_job", sources.HandleJobDeleteByID, projectAuthMiddleware)

	// todo: сделать роуты для обновления проектов, задач

	// Tasks
	tasks := project.Group("/tasks")
	tasks.PUT("/:id_project/create", taskpkg.HandleTaskCreate, projectAuthMiddleware)
	tasks.GET("/:id_project/all", taskpkg.HandleTasksGetList, projectAuthMiddleware)
	tasks.GET("/:id_project/:id_task", taskpkg.HandleTaskGetByID, projectAuthMiddleware)
	tasks.DELETE("/:id_project/:id_task", taskpkg.HandleTaskDeleteByID, projectAuthMiddleware)

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
