package models

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"time"
)

var stormList map[string]*storm.DB
var stormDir = "./monitordb"

func init() {
	stormList = make(map[string]*storm.DB)
	if !utils.Exist(stormDir) {
		_ = os.MkdirAll(stormDir, 0755)
	}
}

type NodeInfoData struct {
	Id           int                            `storm:"id,increment" json:"id"`    //主键,自增长
	Memory       *mem.VirtualMemoryStat         `json:"memory" `                    //内存信息
	NetInterface []net.InterfaceStat            `json:"net_interface"`              //网络硬件信息
	NetIo        []net.IOCountersStat           `json:"net_io"`                     //网格IO使用情况
	CpuUse       *CpuUse                        `json:"cpu_use"`                    //CPU 占用情况
	DiskUse      []*disk.UsageStat              `json:"disk_use"`                   //磁盘使用信息
	DiskIo       map[string]disk.IOCountersStat `json:"disk_io"`                    //硬盘io统计
	CreatedDate  int64                          `storm:"index" json:"created_date"` //数据接收时间
}

type CpuUse struct {
	All  float64   `json:"all"`
	List []float64 `json:"list"`
}

// 表名
type NodeInfoModel struct {
	Table string `json:"table"` //表名
	Order string `json:"order"` //默认排序
	storm.Node
}

func NewNodeInfoModel(nodeAddr string) *NodeInfoModel {
	var db *storm.DB
	var ok bool
	if db, ok = stormList[nodeAddr]; !ok {
		var err error
		db, err = initStormDb(nodeAddr)
		if err != nil {
			return nil
		}
		stormList[nodeAddr] = db
	}

	return &NodeInfoModel{
		Table: "node_info",
		Node:  db.From("node_info"),
		Order: "DESC",
	}
}

func initStormDb(nodeAddr string) (*storm.DB, error) {
	file := fmt.Sprintf("%s.sdb", nodeAddr)
	allFile := fmt.Sprintf("%s/%s", stormDir, file)
	return storm.Open(allFile)
}

// 通过ID拿到记录
func (n *NodeInfoModel) GetById(id int) (*NodeInfoData, error) {
	data := &NodeInfoData{}
	err := n.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetByNodeId 通过节点id得到信息
func (n *NodeInfoModel) GetByNodeId(id int) (*NodeInfoData, error) {
	data := &NodeInfoData{}
	err := n.One("NodeId", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n *NodeInfoModel) SetOrder(ord string) {
	n.Order = ord
}

// 查询条件得到任务数据列表
func (n *NodeInfoModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []NodeInfoData
	count, err := n.Select(where...).Count(new(NodeInfoData))
	if err != nil {
		return nil, err
	}
	query := n.Select(where...)
	if n.Order == "DESC" {
		query.Reverse()
	}

	err = query.Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		List:  list,
		Count: count,
	}, nil
}

// 查询条件得到任务数据列表
func (n *NodeInfoModel) List(page, number int, where ...q.Matcher) ([]NodeInfoData, error) {
	var list []NodeInfoData
	query := n.Select(where...)
	if n.Order == "DESC" {
		query.Reverse()
	}
	err := query.Limit(number).Skip((page - 1) * number).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (n *NodeInfoModel) SaveRange(data *NodeInfoData) error {
	err := n.deleteRange()
	if err != nil {
		return err
	}

	return n.Save(data)
}

// 删除超过一天的记录
func (n *NodeInfoModel) deleteRange() error {
	rangeDay := time.Now().Unix()

	var list []NodeInfoData
	query := n.Select(q.Lt("CreatedDate", rangeDay-(24*3600)))
	num, err := query.Count(new(NodeInfoData))
	if err != nil {
		return err
	}
	if num < 100 {
		return nil
	}

	err = query.Find(&list)
	if err != nil {
		return err
	}
	tx, err := n.Begin(true)
	if err != nil {
		return err
	}
	for _, v := range list {
		err = tx.DeleteStruct(&v)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
