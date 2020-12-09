package middles

import (
	"github.com/clakeboy/golib/components"
	"github.com/gin-gonic/gin"
)

func Redis() gin.HandlerFunc {
	return func(c *gin.Context) {
		rd, _ := components.NewCKRedis()
		c.Set("redis", rd)
		defer rd.Close()
		c.Next()
	}
}
