package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
	"system-monitoring/models"
)

//控制器
type ShellManagerController struct {
	c *gin.Context
}

func NewShellManagerController(c *gin.Context) *ShellManagerController {
	return &ShellManagerController{c: c}
}

// ActionQuery 查询
func (s *ShellManagerController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		Query  []*Condition `json:"query"`
		Page   int          `json:"page"`
		Number int          `json:"number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}
	where := explainQueryCondition(params.Query)
	model := models.NewShellModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ActionSave 添加
func (s *ShellManagerController) ActionSave(args []byte) error {

	return nil
}

// ActionDelete 删除
func (s *ShellManagerController) ActionDelete(args []byte) error {
	return nil
}

func (s *ShellManagerController) ActionFind(args []byte) (*models.ShellData, error) {
	var params struct {
		Id int `json:"id"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewShellModel(nil)
	return model.GetById(params.Id)
}
