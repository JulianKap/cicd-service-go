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
	allProjectsAuthMiddleware := AllProjectsAuthMiddleware()
	projectAuthMiddleware := ProjectAuthMiddleware()

	// Projects
	project := e.Group("/project")
	project.PUT("/create", sources.HandleProjectCreate, allProjectsAuthMiddleware)
	project.GET("/all", sources.HandleProjectsGetList, allProjectsAuthMiddleware)
	project.GET("/:id", sources.HandleProjectGetByID, allProjectsAuthMiddleware)
	project.DELETE("/:id", sources.HandleProjectDeleteByID, allProjectsAuthMiddleware)

	// Jobs
	jobs := project.Group("/jobs")
	jobs.PUT("/create", sources.HandleJobCreate, projectAuthMiddleware)
	jobs.GET("/:id_project/all", sources.HandleJobsGetList, projectAuthMiddleware)
	jobs.GET("/:id_project/:id_job", sources.HandleJobGetByID, projectAuthMiddleware)
	jobs.DELETE("/:id_project/:id_job", sources.HandleJobDeleteByID, projectAuthMiddleware)

	// todo: сделать роуты для обновления проектов, задач

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

func AllProjectsAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			// todo: брать токен из файла конфигурации
			test_current_token := "12345678"
			if token != "Bearer "+test_current_token {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			return next(c)
		}
	}
}

func ProjectAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//token := c.Request().Header.Get("Authorization")
			//projectID := c.Param("id_project")
			//
			// todo: брать токен для конкретного проекта из vault
			//// Получить токен из vault
			//
			//if token != secret.Data["token"].(string) {
			//	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			//}

			// Продолжаем выполнение цепочки обработчиков, если токен действителен
			return next(c)
		}
	}
}
