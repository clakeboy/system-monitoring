package router

import (
	"github.com/gin-gonic/gin"
	"system-monitoring/controllers"
)

func GetController(controllerName string, c *gin.Context) interface{} {
	switch controllerName {
	case "def":
		return controllers.NewDefaultController(c)
	case "login":
		return controllers.NewLoginController(c)
	case "manager":
		return controllers.NewManagerController(c)
	case "node":
		return controllers.NewNodesController(c)
	default:
		return nil
	}
}
