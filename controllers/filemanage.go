package controllers

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
	"system-monitoring/service"
)

// 控制器
type FileManageController struct {
	c *gin.Context
}

func NewFileManageController(c *gin.Context) *FileManageController {
	return &FileManageController{c: c}
}

// ActionQuery 查询
func (f *FileManageController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		Ip     string `json:"ip"`
		Path   string `json:"path"`
		Page   int    `json:"page"`
		Number int    `json:"number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	server, err := service.MainServer.GetNodeServer(params.Ip)
	if err != nil {
		return nil, err
	}

	list, count, err := server.GetRemoteDir(params.Path, params.Page, params.Number)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// ActionGetFile 得到文件内容
func (f *FileManageController) ActionGetFile(args []byte) (string, error) {
	var params struct {
		Ip   string `json:"ip"`
		Path string `json:"path"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return "", err
	}

	server, err := service.MainServer.GetNodeServer(params.Ip)
	if err != nil {
		return "", err
	}

	content, err := server.GetRemoteFile(params.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (f *FileManageController) ActionSaveFile(args []byte) error {
	var params struct {
		Ip      string `json:"ip"`
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	server, err := service.MainServer.GetNodeServer(params.Ip)
	if err != nil {
		return err
	}
	err = server.SaveRemoteFile(params.Path, []byte(params.Content))
	return err
}
