package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
)

//控制器
type FileManageController struct {
	c *gin.Context
}

func NewFileManageController(c *gin.Context) *FileManageController {
	return &FileManageController{c: c}
}

//查询
func (f *FileManageController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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
