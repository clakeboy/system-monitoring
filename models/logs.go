package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"system-monitoring/common"
)

// LogsData 操作日志记录
type LogsData struct {
	Id int `storm:"id,increment" json:"id"` //主键,自增长

}

//表名
type LogsModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewLogsModel(db *storm.DB) *LogsModel {
	if db == nil {
		db = common.BDB
	}

	return &LogsModel{
		Table: "logs",
		Node:  db.From("logs"),
	}
}

//通过ID拿到记录
func (l *LogsModel) GetById(id int) (*LogsData, error) {
	data := &LogsData{}
	err := l.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//查询条件得到任务数据列表
func (l *LogsModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []LogsData
	count, err := l.Select(where...).Count(new(LogsData))
	if err != nil {
		return nil, err
	}
	err = l.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

//查询条件得到任务数据列表
func (l *LogsModel) List(page, number int, where ...q.Matcher) ([]LogsData, error) {
	var list []LogsData
	err := l.Select(where...).Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
