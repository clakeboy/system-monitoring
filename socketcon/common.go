package socketcon

import (
	"bytes"
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
