package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/httputils"
	"github.com/clakeboy/golib/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"system-monitoring/components"
	"system-monitoring/models"
)

// LoginController 登录控制器
type LoginController struct {
	c *gin.Context
}

func NewLoginController(c *gin.Context) *LoginController {
	return &LoginController{c: c}
}

// ActionAuth 验证是否已登录
func (l *LoginController) ActionAuth(args []byte) (*models.ManagerData, error) {
	manager, err := components.AuthUser(l.c)
	if err != nil {
		return nil, err
	}
	manager.Password = ""

	return manager, nil
}

// ActionLogin 登录
func (l *LoginController) ActionLogin(args []byte) (*models.ManagerData, error) {
	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	manager := new(models.ManagerData)

	model := models.NewManagerModel(nil)
	err = model.One("Account", params.Username, manager)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if utils.EncodeMD5(params.Password) != manager.Password {
		return nil, fmt.Errorf("用户或密码错误")
	}
	manager.Password = ""

	cookie := l.c.MustGet("cookie").(*httputils.HttpCookie)
	cookie.Set(components.CookieName, strconv.Itoa(manager.Id), 7*24*3600)

	return manager, nil
}

// ActionLogout 退出登录
func (l *LoginController) ActionLogout(args []byte) error {
	cookie := l.c.MustGet("cookie").(*httputils.HttpCookie)
	cookie.Delete(components.CookieName)
	return nil
}
