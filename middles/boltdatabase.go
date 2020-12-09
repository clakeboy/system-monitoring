package middles

import (
	"github.com/gin-gonic/gin"
	"system-monitoring/common"
)

func BoltDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("bolt", common.BDB)
		c.Next()
	}
}
