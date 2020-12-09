package middles

import (
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
	"system-monitoring/common"
)

func Mongo() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := ckdb.NewDB(common.Conf.MDB.DBName)
		c.Set("mdb", db)
		defer db.Close()
		c.Next()
	}
}
