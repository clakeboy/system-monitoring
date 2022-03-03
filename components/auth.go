package components

import (
	"github.com/clakeboy/golib/httputils"
	"github.com/gin-gonic/gin"
	"strconv"
	"system-monitoring/models"
)

const CookieName = "sys_acc"

func AuthUser(c *gin.Context) (*models.ManagerData, error) {
	cookie := c.MustGet("cookie").(*httputils.HttpCookie)
	acc, err := cookie.Get(CookieName)

	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(acc)
	if err != nil {
		return nil, err
	}

	model := models.NewManagerModel(nil)
	return model.GetById(id)
}
