package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

type AccountGroupData struct {
	Id int `storm:"id,increment" json:"id"` //主键,自增长

}

// AccountGroupModel 用户帐号组
type AccountGroupModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewAccountGroupModel(db *storm.DB) *AccountGroupModel {
	if db == nil {
		db = common.BDB
	}

	return &AccountGroupModel{
		Table: "acc_group",
		Node:  db.From("acc_group"),
	}
}

// GetById 通过ID拿到记录
func (a *AccountGroupModel) GetById(id int) (*AccountGroupData, error) {
	data := &AccountGroupData{}
	err := a.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Query 查询条件得到任务数据列表
func (a *AccountGroupModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []AccountGroupData
	count, err := a.Select(where...).Count(new(AccountGroupData))
	if err != nil {
		return nil, err
	}
	err = a.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// List 查询条件得到任务数据列表
func (a *AccountGroupModel) List(page, number int, where ...q.Matcher) ([]AccountGroupData, error) {
	var list []AccountGroupData
	err := a.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
