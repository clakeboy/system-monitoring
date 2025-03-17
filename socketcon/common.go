package socketcon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"strings"
	"system-monitoring/components"
)

const (
	StatusClose  = "close"
	StatusOpen   = "open"
	StatusActive = "active"
)

// Auth 验证消息内容
type Auth struct {
	Auth string `json:"auth"` //验证密钥
	Name string `json:"name"` //节点名称
}

func (a *Auth) Build() []byte {
	var buf bytes.Buffer
	buf.Write(components.BuildStreamData([]byte(a.Auth)))
	buf.Write(components.BuildStreamData([]byte(a.Name)))
	return buf.Bytes()
}

func (a *Auth) Parse(data []byte) {
	list := components.ParseStreamData(data)
	a.Auth = string(list[0])
	a.Name = string(list[1])
}

// CmdShell 执行命令内容
type CmdShell struct {
	Cmd        string   `json:"cmd"`         //要执行的命令
	Args       []string `json:"args"`        //执行的命令参数
	Dir        string   `json:"dir"`         //执行命令的目录
	AckId      string   `json:"ack_id"`      //执行命令回执编号
	AckContent []byte   `json:"ack_content"` //执行命令回执内容
}

func (c *CmdShell) Build() []byte {
	var buf bytes.Buffer
	buf.Write(components.BuildStreamData([]byte(c.Cmd)))
	buf.Write(components.BuildStreamData([]byte(strings.Join(c.Args, " "))))
	buf.Write(components.BuildStreamData([]byte(c.Dir)))
	buf.Write(components.BuildStreamData([]byte(c.AckId)))
	buf.Write(components.BuildStreamData(c.AckContent))
	return buf.Bytes()
}

func (c *CmdShell) Parse(data []byte) {
	list := components.ParseStreamData(data)
	c.Cmd = string(list[0])
	c.Args = strings.Split(string(list[1]), " ")
	c.Dir = string(list[2])
	c.AckId = string(list[3])
	c.AckContent = list[4]
}

// BuildExec 编译执行命令语句
func (c *CmdShell) BuildExec() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}

// CMDFileInfo 推送文件
type CMDFileInfo struct {
	FileId  int    //文件id
	Path    string //文件替换地址
	FileUri string //文件下载地址
	Error   string //文件处理错误信息
	Message string //文件处理回执信息
}

func (c *CMDFileInfo) Build() []byte {
	var buf bytes.Buffer
	buf.Write(components.BuildStreamData(utils.IntToBytes(c.FileId, 32)))
	buf.Write(components.BuildStreamData([]byte(c.Path)))
	buf.Write(components.BuildStreamData([]byte(c.FileUri)))
	buf.Write(components.BuildStreamData([]byte(c.Error)))
	buf.Write(components.BuildStreamData([]byte(c.Message)))
	return buf.Bytes()
}

func (c *CMDFileInfo) Parse(data []byte) {
	list := components.ParseStreamData(data)
	c.FileId = utils.BytesToInt(list[0])
	c.Path = string(list[1])
	c.FileUri = string(list[2])
	c.Error = string(list[3])
	c.Message = string(list[4])
}

// CMDFileTrans 文件传送
type CMDFileTrans struct {
	Name    string //文件名
	Path    string //文件替换地址
	Content []byte //文件内容
}

func (c *CMDFileTrans) Build() {

}

func (c *CMDFileTrans) Parse() {

}

type CMDDirType int

const (
	DirList CMDDirType = iota + 1
	DirContent
	DirSaveFile
)

func (c CMDDirType) String() string {
	switch c {
	case DirList:
		return "DirList"
	case DirContent:
		return "DirContent"
	case DirSaveFile:
		return "DirSaveFile"
	default:
		return ""
	}
}

// CMDDir 文件目录
type CMDDir struct {
	Sn        string        `json:"sn"`         //调用序列号
	Path      string        `json:"path"`       //目录地址
	Page      int           `json:"page"`       //列表当前页
	PageCount int           `json:"page_count"` //总页数
	Count     int           `json:"count"`      //总列表数
	Number    int           `json:"number"`     //列表显示数量
	Type      CMDDirType    `json:"type"`       //数据类型
	Error     string        `json:"error"`      //错误信息
	Content   []byte        `json:"content"`    //文件内容
	List      []*CMDDirList `json:"list"`
}

type CMDDirList struct {
	Name         string `json:"name"`          //文件名
	IsDir        bool   `json:"is_dir"`        //是否目录
	Size         int64  `json:"size"`          //文件大小
	Mode         string `json:"mode"`          //模式
	ModifiedDate int64  `json:"modified_date"` //修改时间
}

func (c *CMDDir) Build() []byte {
	var buf bytes.Buffer
	data, err := json.Marshal(c)
	if err == nil {
		buf.Write(components.BuildStreamData(data))
	}
	return buf.Bytes()
}

func (c *CMDDir) Parse(data []byte) error {
	list := components.ParseStreamData(data)
	err := json.Unmarshal(list[0], c)
	return err
}
