package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

type ManagerData struct {
	Id          int    `storm:"id,increment" json:"id"` //主键,自增长
	Account     string `storm:"index" json:"account"`   //登录名
	Password    string `json:"password"`                //密码
	Phone       string `json:"phone"`                   //电话
	Name        string `json:"name"`                    //名称
	CreatedDate int64  `json:"created_date"`            //创建时间
	CreatedBy   string `json:"created_by"`              //创建人
}

// ManagerModel 表名
type ManagerModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewManagerModel(db *storm.DB) *ManagerModel {
	if db == nil {
		db = common.BDB
	}

	return &ManagerModel{
		Table: "manager",
		Node:  db.From("manager"),
	}
}

// GetById 通过ID拿到记录
func (m *ManagerModel) GetById(id int) (*ManagerData, error) {
	data := &ManagerData{}
	err := m.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
func (m *ManagerModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []ManagerData
	count, err := m.Select(where...).Count(new(ManagerData))
	if err != nil {
		return nil, err
	}
	err = m.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// List 查询条件得到任务数据列表
func (m *ManagerModel) List(page, number int, where ...q.Matcher) ([]ManagerData, error) {
	var list []ManagerData
	//m.Select().OrderBy("user")
	err := m.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
