package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

const (
	NodeStatusOnline = iota + 1
	NodeStatusOffline
)

// NodeData 服务节点数据
type NodeData struct {
	Id             int    `storm:"id,increment" json:"id"` //主键,自增长
	Name           string `storm:"index" json:"name"`      //节点名称
	Ip             string `storm:"index" json:"ip"`        //节点IP地址
	Status         int    `json:"status"`                  //节点状态
	LastOnlineDate int64  `json:"last_online_date"`        //最后一次在线时间
	CreateDate     int64  `json:"create_date"`             //第一次创建时间
}

// NodeModel 表名
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

// GetById 通过ID拿到记录
func (n *NodeModel) GetById(id int) (*NodeData, error) {
	data := &NodeData{}
	err := n.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetByIp 通过IP拿到记录
func (n *NodeModel) GetByIp(ipAddr string) (*NodeData, error) {
	data := &NodeData{}
	err := n.One("Ip", ipAddr, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
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

// List 查询条件得到任务数据列表
func (n *NodeModel) List(page, number int, where ...q.Matcher) ([]NodeData, error) {
	var list []NodeData
	err := n.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
