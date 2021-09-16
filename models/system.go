package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

// SystemData 系统信息
type SystemData struct {
	Id     int    `storm:"id,increment" json:"id"` //主键,自增长
	NodeId int    `storm:"index" json:"node_id"`   //服务节点ID
	Name   string `json:"system_name"`             //系统名称
	Path   string `json:"path"`                    //系统路径

}

// SystemModel 更新系统
type SystemModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewSystemModel(db *storm.DB) *SystemModel {
	if db == nil {
		db = common.BDB
	}

	return &SystemModel{
		Table: "system",
		Node:  db.From("system"),
	}
}

// GetById 通过ID拿到记录
func (s *SystemModel) GetById(id int) (*SystemData, error) {
	data := &SystemData{}
	err := s.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
func (s *SystemModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []SystemData
	count, err := s.Select(where...).Count(new(SystemData))
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
func (s *SystemModel) List(page, number int, where ...q.Matcher) ([]SystemData, error) {
	var list []SystemData
	err := s.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
