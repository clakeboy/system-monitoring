package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

//服务节点数据
type NodeData struct {
	Id int `storm:"id,increment" json:"id"` //主键,自增长
}

//表名
type NodeModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewNodeModel(db *storm.DB) *NodeModel {
	if db == nil {
		db = common.BDB
	}

	return &NodeModel{
		Table: "node",
		Node:  db.From("node"),
	}
}

//通过ID拿到记录
func (n *NodeModel) GetById(id int) (*NodeData, error) {
	data := &NodeData{}
	err := n.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//查询条件得到任务数据列表
func (n *NodeModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []NodeData
	count, err := n.Select(where...).Count(new(NodeData))
	if err != nil {
		return nil, err
	}
	err = n.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

//查询条件得到任务数据列表
func (n *NodeModel) List(page, number int, where ...q.Matcher) ([]NodeData, error) {
	var list []NodeData
	err := n.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
