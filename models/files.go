package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

//上传的服务文件
type FileData struct {
	Id           int    `storm:"id,increment" json:"id"`  //主键,自增长
	ServiceId    int    `storm:"index" json:"service_id"` //上传的服务id
	Name         string `storm:"index" json:"name"`       //文件名
	OrgName      string `json:"org_name"`                 //原始文件名
	Path         string `json:"path"`                     //文件所在目录
	Size         int64  `json:"size"`                     //文件大小
	Type         string `json:"type"`                     //文件类型
	CreatedDate  int64  `json:"created_date"`             //创建时间
	CreatedBy    string `json:"created_by"`               //创建人
	PushResult   string `json:"push_result"`              //推送结果
	PushLastTime int64  `json:"push_last_time"`           //最后推送时间
}

//表名
type FileModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewFileModel(db *storm.DB) *FileModel {
	if db == nil {
		db = common.BDB
	}

	return &FileModel{
		Table: "file_list",
		Node:  db.From("file_list"),
	}
}

//通过ID拿到记录
func (f *FileModel) GetById(id int) (*FileData, error) {
	data := &FileData{}
	err := f.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//查询条件得到任务数据列表
func (f *FileModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []FileData
	count, err := f.Select(where...).Count(new(FileData))
	if err != nil {
		return nil, err
	}
	err = f.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

//查询条件得到任务数据列表
func (f *FileModel) List(page, number int, where ...q.Matcher) ([]FileData, error) {
	var list []FileData
	err := f.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
