package socketcon

import (
	"bytes"
	"fmt"
	"github.com/clakeboy/golib/components"
	"os/exec"
	components2 "system-monitoring/components"
)

// NodeClient 节点控制
type NodeClient struct {
	conn   *components.TCPConnect
	log    *components.SysLog
	status string
	events map[string]func(evt *components.TCPConnEvent)
}

// NewNodeClient 创建一个新的节点客户端
func NewNodeClient() *NodeClient {
	return &NodeClient{
		log:    components.NewSysLog("node_client_"),
		status: StatusOpen,
		events: make(map[string]func(evt *components.TCPConnEvent)),
	}
}

// OnConnected 连接完成事件
func (n *NodeClient) OnConnected(e *components.TCPConnEvent) {
	n.conn = e.Conn
	if evt, ok := n.events["connected"]; ok {
		e.Data = n
		evt(e)
	}
}

// OnDisconnected 关闭连接
func (n *NodeClient) OnDisconnected(e *components.TCPConnEvent) {
	if evt, ok := n.events["disconnect"]; ok {
		evt(e)
	}
}

// OnRecv 接收数据
func (n *NodeClient) OnRecv(evt *components.TCPConnEvent) {
	fmt.Println("server recv:", evt)
	data := evt.Data.([]byte)
	if len(data) <= 0 {
		return
	}
	cmd := components2.NewMainStream()
	err := cmd.Parse(evt.Data.([]byte))
	if err != nil {
		n.log.Error(fmt.Errorf("receive data error: %v", err))
		n.conn.Close()
		return
	}

	n.execCommand(cmd)
}

// OnWritten 写入数据后
func (n *NodeClient) OnWritten(evt *components.TCPConnEvent) {
	//p.conn.Close()
}

// OnError 错误事件
func (n *NodeClient) OnError(evt *components.TCPConnEvent) {
	fmt.Println(evt.Data)
}

//执行命令
func (n *NodeClient) execCommand(cmd *components2.MainStream) {
	fmt.Println("exec command:", cmd.Command)
	switch cmd.Command {
	case components2.CMDPing:
		fmt.Println("server:", n.conn.RemoteAddr(), "-> ping")
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDClose:
		n.conn.Close()
		n.status = StatusClose
	case components2.CMDAuthCode:
		n.status = StatusActive
		fmt.Println("login done")
		if evt, ok := n.events["login"]; ok {
			evt(&components.TCPConnEvent{
				Data: cmd,
			})
		}
	case components2.CMDShell:
		shell := new(CmdShell)
		shell.Parse(cmd.Content)
		n.execShell(shell)
	}
}

// On 外部绑定事件
func (n *NodeClient) On(evtName string, evt func(evt *components.TCPConnEvent)) {
	n.events[evtName] = evt
}

//执行shell命令并返回结果
func (n *NodeClient) execShell(cmd *CmdShell) {
	shell := exec.Command(cmd.Cmd, cmd.Args...)
	shell.Dir = cmd.Dir
	var buf bytes.Buffer
	shell.Stdout = &buf
	shell.Stderr = &buf
	err := shell.Run()
	if err != nil {
		buf.WriteString(fmt.Sprintf("[ERROR] %s", err))
	}

	shellData := &CmdShell{
		AckId:      cmd.AckId,
		AckContent: buf.Bytes(),
	}

	ackCmd := components2.NewMainStream()
	ackCmd.Command = components2.CMDShell
	ackCmd.Content = shellData.Build()
	n.conn.WriteData(ackCmd.Build())
}
