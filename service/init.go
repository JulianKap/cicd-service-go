package service

import (
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
	//handler := schedule.NewHandler()
	//api := e.Group("/api")
	//
	//api.GET("/version", schedule.HandleVersion)
	//api.GET("/master", schedule.HandleMaster)
	//api.GET("/slave", schedule.HandleSlave)
	//
	////
	//api.POST("/create", schedule.HandleScheduleByID)
	//api.POST("/delete", schedule.HandleScheduleNew)
	//
	//api.PUT("/namespace/", schedule.HandleScheduleUpdateByID)
	//api.DELETE("/:id", schedule.HandleScheduleDeleteByID)
	//
	//api.POST("/project", schedule.HandleScheduleByID)
	//api.POST("/job", schedule.HandleScheduleNew)
	//
	//// Namespace
	//
	//// Job
	//api.POST("/:id", schedule.HandleScheduleUpdateByID)
	//// Удалить job
	//api.DELETE("/:id", schedule.HandleScheduleDeleteByID)

	e.GET("/hc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
}
