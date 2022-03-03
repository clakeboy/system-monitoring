package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
	"system-monitoring/models"
	"system-monitoring/service"
	"time"
)

//控制器
type ServiceController struct {
	c *gin.Context
}

func NewServiceController(c *gin.Context) *ServiceController {
	return &ServiceController{c: c}
}

// ActionQuery 查询
func (s *ServiceController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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
	model := models.NewServiceModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ActionSave 保存
func (s *ServiceController) ActionSave(args []byte) error {
	var params struct {
		Data *models.ServiceData `json:"data"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	saveData := params.Data

	model := models.NewServiceModel(nil)

	if saveData.Id == 0 {
		saveData.CreatedDate = time.Now().Unix()
		return model.Save(saveData)
	}

	orgData, err := model.GetById(saveData.Id)
	if err != nil {
		return err
	}

	orgData.Name = saveData.Name
	orgData.Type = saveData.Type
	orgData.NodeName = saveData.NodeName
	orgData.NodeId = saveData.NodeId
	orgData.Directory = saveData.Directory
	orgData.ModifiedDate = time.Now().Unix()

	return model.Update(orgData)
}

// ActionFind 查询数据
func (s *ServiceController) ActionFind(args []byte) (*models.ServiceData, error) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewServiceModel(nil)
	data, err := model.GetById(params.Id)

	return data, err
}

// ActionDelete 删除
func (s *ServiceController) ActionDelete(args []byte) error {
	return nil
}

// ActionExec 执行命令
func (s *ServiceController) ActionExec(args []byte) (*models.ShellData, error) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewServiceModel(nil)
	data, err := model.GetById(params.Id)
	if err != nil {
		return nil, err
	}

	nodeModel := models.NewNodeModel(nil)
	node, err := nodeModel.GetById(data.NodeId)
	if err != nil {
		return nil, err
	}

	ok := service.MainServer.CheckIp(node.Ip)
	if !ok {
		return nil, fmt.Errorf("can not execute update ,node server '%s:%s' offline", node.Name, node.Ip)
	}
	shell := &models.ShellData{
		NodeId:      node.Id,
		NodeName:    node.Name,
		Cmd:         "git",
		Args:        []string{"pull"},
		Dir:         data.Directory,
		Status:      0,
		ExecBy:      "",
		ExecContent: "",
		CreateDate:  time.Now().Unix(),
	}
	shellModel := models.NewShellModel(nil)
	err = shellModel.Save(shell)
	if err != nil {
		return nil, err
	}
	err = service.MainServer.ExecShell(node.Ip, shell)
	return shell, err
}

// ActionExecResult 查询命令结果
func (s *ServiceController) ActionExecResult(args []byte) (*models.ShellData, error) {
	var params struct {
		ShellId int `json:"shell_id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewShellModel(nil)
	shell, err := model.GetById(params.ShellId)
	if err != nil {
		return nil, err
	}

	return shell, err
}
