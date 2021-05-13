package socketcon

import (
	"bytes"
	"system-monitoring/components"
)

const (
	StatusClose  = "close"
	StatusOpen   = "open"
	StatusActive = "active"
)

//验证消息内容
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
