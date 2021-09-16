package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/gin-gonic/gin"
	"system-monitoring/models"
	"time"
)

// ManagerController 管理人员控制器
type ManagerController struct {
	c *gin.Context
}

func NewManagerController(c *gin.Context) *ManagerController {
	return &ManagerController{c: c}
}

// ActionQuery 查询
func (m *ManagerController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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
	model := models.NewManagerModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ActionSave 保存
func (m *ManagerController) ActionSave(args []byte) error {
	var params struct {
		Data *models.ManagerData `json:"data"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	saveData := params.Data

	model := models.NewManagerModel(nil)

	if saveData.Id == 0 {
		saveData.CreatedDate = time.Now().Unix()
		saveData.Password = utils.EncodeMD5(saveData.Password)
		return model.Save(saveData)
	}

	orgData, err := model.GetById(saveData.Id)
	if err != nil {
		return err
	}

	orgData.Name = saveData.Name
	orgData.Phone = saveData.Phone
	if saveData.Password != "" {
		orgData.Password = utils.EncodeMD5(saveData.Password)
	}

	return model.Update(orgData)
}

// ActionFind 查找用户
func (m *ManagerController) ActionFind(args []byte) (*models.ManagerData, error) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewManagerModel(nil)
	data, err := model.GetById(params.Id)
	if err != nil {
		return nil, err
	}
	data.Password = ""
	return data, err
}

// ActionDelete 删除
func (m *ManagerController) ActionDelete(args []byte) error {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	model := models.NewManagerModel(nil)
	err = model.DeleteStruct(&models.ManagerData{
		Id: params.Id,
	})

	return err
}
