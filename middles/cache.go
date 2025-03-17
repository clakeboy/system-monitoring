package middles

import (
	"github.com/gin-gonic/gin"
	"system-monitoring/common"
)

func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", common.MemCache)
		c.Next()
	}
}
