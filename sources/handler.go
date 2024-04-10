package sources

import (
	"cicd-service-go/init/db"
	"cicd-service-go/manager"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleProjectCreate(ctx echo.Context) (err error) {
	var project Project

	project = Project{
		ID:          "1",
		APIKey:      "1342342342342",
		ProjectName: "test",
	}

	if !manager.MemberInfo.Master {
		log.Info("HandleProjectCreate info #0: not master")
		return ctx.JSON(http.StatusBadRequest, Response{Message: "I am SLAVE! Slave not support create projects"})
	}

	err = project.createProject(db.InstanceETCD)
	if err != nil {
		log.Error("HandleProjectCreate[0]: Ошибка: ", err)
		return ctx.JSON(http.StatusBadRequest, Response{Message: "Error"})
	}

	return ctx.JSON(http.StatusOK, Response{Message: "Project created"})
}

func HandleJobCreate(ctx echo.Context) (err error) {

	return ctx.JSON(http.StatusOK, Response{Message: "Job created"})
}
