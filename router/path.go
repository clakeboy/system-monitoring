package router

import (
	"github.com/gin-gonic/gin"
	"system-monitoring/controllers"
)

func GetController(controllerName string, c *gin.Context) interface{} {
	switch controllerName {
	case "def":
		return controllers.NewDefaultController(c)
	default:
		return nil
	}
}
