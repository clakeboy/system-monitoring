package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

type ShellData struct {
	Id          int      `storm:"id,increment" json:"id"` //主键,自增长
	NodeId      int      `storm:"index" json:"node_id"`   //服务节点ID
	NodeName    string   `json:"node_name"`               //服务节点名称
	Cmd         string   `json:"cmd"`                     //执行的命令
	Args        []string `json:"args"`                    //命令参数
	Dir         string   `json:"dir"`                     //执行命令目录
	Status      int      `json:"status"`                  //执行状态
	ExecBy      string   `json:"exec_by"`                 //执行人
	ExecContent string   `json:"exec_content"`            //执行内容
	CreateDate  int64    `json:"create_date"`             //执行时间
}

// ShellModel 执行shell命令日志表
type ShellModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewShellModel(db *storm.DB) *ShellModel {
	if db == nil {
		db = common.BDB
	}

	return &ShellModel{
		Table: "shell",
		Node:  db.From("shell"),
	}
}

// GetById 通过ID拿到记录
func (s *ShellModel) GetById(id int) (*ShellData, error) {
	data := &ShellData{}
	err := s.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
func (s *ShellModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []ShellData
	count, err := s.Select(where...).Count(new(ShellData))
	if err != nil {
		return nil, err
	}
	err = s.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// List 查询条件得到任务数据列表
func (s *ShellModel) List(page, number int, where ...q.Matcher) ([]ShellData, error) {
	var list []ShellData
	err := s.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
