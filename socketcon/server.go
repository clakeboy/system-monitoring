package socketcon

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"system-monitoring/common"
	components2 "system-monitoring/components"
	"time"
)

//节点控制
type NodeServer struct {
	conn            *components.TCPConnect
	log             *components.SysLog
	status          string
	name            string
	EventDisconnect func(evt *components.TCPConnEvent)
}

//创建一个新的主服务
func NewNodeServer() *NodeServer {
	return &NodeServer{
		log:    components.NewSysLog("node_server_"),
		status: StatusClose,
	}
}

//连接完成事件
func (n *NodeServer) OnConnected(evt *components.TCPConnEvent) {
	n.conn = evt.Conn
	n.Ping()
}

//关闭连接
func (n *NodeServer) OnDisconnected(evt *components.TCPConnEvent) {
	n.status = "close"
}

//接收数据
func (n *NodeServer) OnRecv(evt *components.TCPConnEvent) {
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
func (n *NodeServer) OnWritten(evt *components.TCPConnEvent) {
	//p.conn.Close()
	cmd := components2.NewMainStream()
	err := cmd.Parse(evt.Data.([]byte))
	if err != nil {
		return
	}
	if cmd.Command == components2.CMDClose {
		n.conn.Close()
	}
}

//错误事件
func (n *NodeServer) OnError(evt *components.TCPConnEvent) {
	fmt.Println("error:", evt.Data)
}

//外部绑定事件
func (n *NodeServer) On(evtName string, evt func(evt *components.TCPConnEvent)) {
	if evtName == "disconnect" {
		n.EventDisconnect = evt
	}
}

//节点名称
func (n *NodeServer) Name() string {
	return n.name
}

//节点状态
func (n *NodeServer) Status() int {
	return n.conn.Status()
}

//周期性PING
func (n *NodeServer) Ping() {
	go func() {
		tk := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-tk.C:
				if n.conn.Status() == components.TCPStatusDisconnected {
					tk.Stop()
					return
				}
				fmt.Println("ping ->", n.name, n.conn.RemoteAddr())
				cmd := components2.NewMainStream()
				cmd.Command = components2.CMDPing
				n.conn.WriteData(cmd.Build())
			}
		}
	}()
}

//执行命令
func (n *NodeServer) execCommand(cmd *components2.MainStream) {
	fmt.Println("exec command:", cmd.Command)
	switch cmd.Command {
	case components2.CMDPong:
		fmt.Println(n.name, n.conn.RemoteAddr(), "<- pong")
	case components2.CMDPing:
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDAuth:
		auth := new(Auth)
		auth.Parse(cmd.Content)
		fmt.Println("login:", auth.Name)
		if auth.Auth == common.Conf.Server.AuthPass {
			n.name = auth.Name
			ackCmd := components2.NewMainStream()
			ackCmd.Command = components2.CMDAuthCode
			n.conn.WriteData(ackCmd.Build())
		} else {
			ackCmd := components2.NewMainStream()
			ackCmd.Command = components2.CMDClose
			n.conn.WriteData(ackCmd.Build())
		}
	}
}
