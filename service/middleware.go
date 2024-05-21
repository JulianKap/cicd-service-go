package service

import (
	"cicd-service-go/init/secrets"
	"cicd-service-go/sources"
	"cicd-service-go/utility"
	"cicd-service-go/vault"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// allProjectsAuthMiddleware авторизации для управления всеми проектами
func allProjectsAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := utility.RemovePrefixAuthBearer(c.Request().Header.Get("Authorization"))

			// todo: брать токен из файла конфигурации
			if token != "12345678" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			return next(c)
		}
	}
}

// projectAuthMiddleware авторизация для работы с проектом
func projectAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := utility.RemovePrefixAuthBearer(c.Request().Header.Get("Authorization"))
			if token == "" {
				log.Error("projectAuthMiddleware #0: Not found token authorized")
				return echo.NewHTTPError(http.StatusBadRequest, "Not found token authorized")
			}

			projectID, err := strconv.Atoi(c.Param("id_project"))
			if err != nil {
				log.Error("projectAuthMiddleware #1: ", err)
				return echo.NewHTTPError(http.StatusBadRequest, "Bad project_id")
			}

			p := sources.Project{ID: projectID}
			codeValPrj, _ := sources.ValidateProject(&p)
			if codeValPrj != http.StatusOK {
				log.Error("projectAuthMiddleware #2: bad project_id")
				return echo.NewHTTPError(http.StatusBadRequest, "Bad project_id")
			}

			ok, err := checkTokenVault(secrets.InstanceVault, &p, &vault.Token{Token: token})
			if err != nil {
				log.Error("projectAuthMiddleware #3: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
			}

			if !ok {
				log.Info("projectAuthMiddleware #4: invalid token=", token)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			return next(c)
		}
	}
}
