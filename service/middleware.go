package service

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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
			//
			//projectID, err := strconv.Atoi(c.Param("id_project"))
			//if err != nil {
			//	log.Error("ProjectAuthMiddleware #0: ", err)
			//	return echo.NewHTTPError(http.StatusBadRequest, "Bad project_id")
			//}
			//
			//ok, err := СheckTokenVault(secrets.InstanceVault, &sources.Project{ID: projectID}, &vault.Token{Token: token})
			//if err != nil {
			//	log.Error("ProjectAuthMiddleware #1: ", err)
			//	return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
			//}
			//
			//if !ok {
			//	log.Info("ProjectAuthMiddleware #2: invalid token=", token)
			//	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			//}

			return next(c)
		}
	}
}
