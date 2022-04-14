package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/gin-gonic/gin"
	"system-monitoring/models"
	"system-monitoring/service"
)

// NodesController 节点控制器
type NodesController struct {
	c *gin.Context
}

func NewNodesController(c *gin.Context) *NodesController {
	return &NodesController{c: c}
}

// ActionQuery 查询
func (n *NodesController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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
	model := models.NewNodeModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *NodesController) ActionDelete(args []byte) error {
	var params struct {
		Id int `json:"id"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	serviceModel := models.NewServiceModel(nil)
	count, err := serviceModel.Select(q.Eq("NodeId", params.Id)).Count(new(models.ServiceData))
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("节点下有服务, 不能删除该节点服务器")
	}

	model := models.NewNodeModel(nil)
	return model.DeleteStruct(&models.NodeData{
		Id: params.Id,
	})
}

func (n *NodesController) ActionCheckOnline(args []byte) (bool, error) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return false, err
	}

	model := models.NewNodeModel(nil)
	data, err := model.GetById(params.Id)
	if err != nil {
		return false, err
	}

	ok := service.MainServer.CheckIp(data.Ip)

	data.Status = utils.YN(ok, 1, 2).(int)
	err = model.Save(data)

	return ok, err
}
