package socketcon

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	components2 "system-monitoring/components"
)

//节点控制
type NodeClient struct {
	conn            *components.TCPConnect
	log             *components.SysLog
	status          string
	EventDisconnect func(evt *components.TCPConnEvent)
}

//创建一个新的节点客户端
func NewNodeClient() *NodeClient {
	return &NodeClient{
		log:    components.NewSysLog("node_client_"),
		status: StatusOpen,
	}
}

//连接完成事件
func (n *NodeClient) OnConnected(evt *components.TCPConnEvent) {
	n.conn = evt.Conn
}

//关闭连接
func (n *NodeClient) OnDisconnected(evt *components.TCPConnEvent) {
	if n.EventDisconnect != nil {
		n.EventDisconnect(evt)
	}
}

//接收数据
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

//写入数据后
func (n *NodeClient) OnWritten(evt *components.TCPConnEvent) {
	//p.conn.Close()
}

//错误事件
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
	}
}

//外部绑定事件
func (n *NodeClient) On(evtName string, evt func(evt *components.TCPConnEvent)) {
	if evtName == "disconnect" {
		n.EventDisconnect = evt
	}
}
