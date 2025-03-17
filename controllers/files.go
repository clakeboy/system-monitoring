package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"system-monitoring/common"
	"system-monitoring/components"
	"system-monitoring/models"
	"system-monitoring/service"
	"time"
)

//控制器
type FilesController struct {
	c *gin.Context
}

func NewFilesController(c *gin.Context) *FilesController {
	return &FilesController{c: c}
}

//查询
func (f *FilesController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
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
	model := models.NewFileModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	return res, nil
}

//删除文件
func (f *FilesController) ActionDelete(args []byte) error {
	return nil
}

//上传文件
func (f *FilesController) ActionUpload() error {
	user, err := components.AuthUser(f.c)
	if err != nil {
		return err
	}
	servId := f.c.PostForm("service_id")
	if servId == "" {
		return fmt.Errorf("service id can not be empty")
	}
	serviceId, err := strconv.Atoi(servId)
	if err != nil {
		return err
	}
	nowTime := time.Now()
	savePath := fmt.Sprintf("%s/%d/%s", common.Conf.Server.FileDir, serviceId, nowTime.Format("20060102"))
	if !utils.Exist(savePath) {
		err = os.MkdirAll(savePath, 0755)
		if err != nil {
			return err
		}
	}
	fh, err := f.c.FormFile("service_file")
	if err != nil {
		return err
	}

	orgName := f.c.PostForm("file_name")
	extStr := filepath.Ext(orgName)
	saveName := fmt.Sprintf("%s-%s%s", strings.ReplaceAll(orgName, extStr, ""), nowTime.Format("150405"), extStr)
	data := &models.FileData{
		Id:          0,
		ServiceId:   serviceId,
		Name:        saveName,
		OrgName:     orgName,
		Path:        savePath,
		Size:        fh.Size,
		Type:        path.Ext(fh.Filename),
		CreatedDate: nowTime.Unix(),
		CreatedBy:   user.Name,
	}
	model := models.NewFileModel(nil)
	err = model.Save(data)
	if err != nil {
		return err
	}

	err = f.c.SaveUploadedFile(fh, fmt.Sprintf("%s/%s", savePath, saveName))
	return err
}

// ActionGet 得到下载文件
func (f *FilesController) ActionGet() {
	fileId := f.c.Query("fid")
	fId, err := strconv.Atoi(fileId)
	if err != nil {
		_ = f.c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid %v", err))
		return
	}

	model := models.NewFileModel(nil)
	fileData, err := model.GetById(fId)
	if err != nil {
		_ = f.c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid %v", err))
		return
	}

	f.c.FileAttachment(fmt.Sprintf("%s/%s", fileData.Path, fileData.Name), fileData.OrgName)
}

// ActionPushFile 推送文件到远程服务
func (f *FilesController) ActionPushFile(args []byte) error {
	var params struct {
		Id        int `json:"id"`
		ServiceId int `json:"service_id"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	fileModel := models.NewFileModel(nil)
	fileData, err := fileModel.GetById(params.Id)
	if err != nil {
		return err
	}
	fileData.PushLastTime = time.Now().Unix()
	fileData.PushResult = ""
	fileModel.Update(fileData)
	serviceModel := models.NewServiceModel(nil)
	serviceData, err := serviceModel.GetById(params.ServiceId)
	if err != nil {
		return err
	}
	nodeModel := models.NewNodeModel(nil)
	node, err := nodeModel.GetById(serviceData.NodeId)
	if err != nil {
		return err
	}
	server, err := service.MainServer.GetNodeServer(node.Ip)
	if err != nil {
		return err
	}
	server.PushFile(fileData, serviceData)

	return nil
}
