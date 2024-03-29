package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

// ServiceData 服务数据
type ServiceData struct {
	Id             int               `storm:"id,increment" json:"id"` //主键,自增长
	NodeId         int               `storm:"index" json:"node_id"`   //服务节点id
	NodeName       string            `json:"node_name"`               //服务节点id
	Name           string            `storm:"index" json:"name"`      //服务名称
	Type           string            `storm:"index" json:"type"`      //服务类型, git,golang
	Directory      string            `json:"directory"`               //服务器目录地址
	Command        string            `json:"command"`                 //执行的命令
	StopCommand    string            `json:"stop_command"`            //停止命令
	RestartCommand string            `json:"restart_command"`         //重启命令
	CommandList    []*ServiceCommand `json:"command_list"`            //命令列表
	CurrentFileId  int               `json:"current_file_id"`         //当前文件id
	FileName       string            `json:"current_file_name"`       //golang服务文件名
	CreatedDate    int64             `json:"created_date"`            //创建时间
	CreatedBy      string            `json:"created_by"`              //创建人
	ModifiedDate   int64             `json:"modified_date"`           //最后修改时间
	ModifiedBy     string            `json:"modified_by"`             //最后修改人
}

// ServiceCommand 服务自定义命令
type ServiceCommand struct {
	Name    string `json:"name"`    //命令名称
	Command string `json:"command"` //执行命令
}

// ServiceModel 服务目录
type ServiceModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewServiceModel(db *storm.DB) *ServiceModel {
	if db == nil {
		db = common.BDB
	}

	return &ServiceModel{
		Table: "service",
		Node:  db.From("service"),
	}
}

// GetById 通过ID拿到记录
func (s *ServiceModel) GetById(id int) (*ServiceData, error) {
	data := &ServiceData{}
	err := s.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
func (s *ServiceModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []ServiceData
	count, err := s.Select(where...).Count(new(ServiceData))
	if err != nil {
		return nil, err
	}
	err = s.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// List 查询条件得到任务数据列表
func (s *ServiceModel) List(page, number int, where ...q.Matcher) ([]ServiceData, error) {
	var list []ServiceData
	err := s.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
