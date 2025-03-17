package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
)

// DefaultController 控制器
type DefaultController struct {
	c *gin.Context
}

func NewDefaultController(c *gin.Context) *DefaultController {
	return &DefaultController{c: c}
}

// ActionConnect 查询
func (d *DefaultController) ActionConnect(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		Server int    `json:"server"`
		Auth   string `json:"auth"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
