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
	case "service":
		return controllers.NewServiceController(c)
	case "shell":
		return controllers.NewShellManagerController(c)
	default:
		return nil
	}
}

//func ApplyRouter(g *gin.Engine) {
//	g.OPTIONS("*action", func(c *gin.Context) {
//		components.Cross(c, h.isCross, c.Request.Header.Get("Origin"))
//	})
//
//	//POST服务接收
//	g.POST("/serv/:controller/:action", func(c *gin.Context) {
//		components.Cross(c, h.isCross, c.Request.Header.Get("Origin"))
//		controller := GetController(c.Param("controller"), c)
//		components.CallAction(controller, c)
//	})
//	//GET服务
//	g.GET("/serv/:controller/:action", func(c *gin.Context) {
//		controller := GetController(c.Param("controller"), c)
//		components.CallActionGet(controller, c)
//	})
//}
