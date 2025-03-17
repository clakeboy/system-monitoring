package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
)

//控制器
type AccountGroupController struct {
	c *gin.Context
}

func NewAccountGroupController(c *gin.Context) *AccountGroupController {
	return &AccountGroupController{c: c}
}

// ActionQuery 查询
func (a *AccountGroupController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		Page   int `json:"page"`
		Number int `json:"number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ActionSave 保存数据
func (a *AccountGroupController) ActionSave(args []byte) error {

	return nil
}

// ActionDelete 删除
func (a *AccountGroupController) ActionDelete(args []byte) error {
	return nil
}
