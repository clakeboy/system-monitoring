package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
)

//控制器
type DefaultController struct {
	c *gin.Context
}

func NewDefaultController(c *gin.Context) *DefaultController {
	return &DefaultController{c: c}
}

//查询
func (d *DefaultController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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

//添加
func (d *DefaultController) ActionInsert(args []byte) error {

	return nil
}

//删除
func (d *DefaultController) ActionDelete(args []byte) error {
	return nil
}

//修改
func (d *DefaultController) ActionUpdate(args []byte) error {
	return nil
}
