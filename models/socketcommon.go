package models

import (
	"encoding/json"
	"github.com/clakeboy/golib/utils"
)

const (
	// SocketEventPty Socket事件列表
	SocketEventPty = "pty"
	// ErrorCode 处理方法代码列表
	ErrorCode    = "error"
	BeginCode    = "start"
	ProcessCode  = "process"
	CompleteCode = "complete"
	CancelCode   = "cancel"
)

const (
	PtyExec     = "exec"
	PtyStart    = "start"
	PtyContinue = "Continue"
	PtyClose    = "close"
)

type PtyMessage struct {
	*utils.JsonParse `json:"-"`
	Evt              string          `json:"evt"`  //pty事件
	Data             json.RawMessage `json:"data"` //数据
	Cmd              string          `json:"cmd"`
}

func NewPtyMessage() *PtyMessage {
	py := &PtyMessage{}
	py.JsonParse = utils.NewJsonParse(py)
	return py
}

// SocketReturn SOCKET 返回数据
type SocketReturn struct {
	*utils.JsonParse `json:"-"`
	Code             string      `json:"code"`    //错误类型
	Message          string      `json:"message"` //错误说明
	Data             interface{} `json:"data"`    //传送数据
}

func NewSocketResult(code, msg string, data interface{}) *SocketReturn {
	js := &SocketReturn{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	js.JsonParse = utils.NewJsonParse(js)
	return js
}
