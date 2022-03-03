package socketcon

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"strings"
	"system-monitoring/common"
	components2 "system-monitoring/components"
	"time"
)

// NodeServerEvent 客户端事件
type NodeServerEvent struct {
	Client *NodeServer //客户端
	Data   interface{} //事件数据
}

// NodeServer 节点控制
type NodeServer struct {
	conn   *components.TCPConnect
	log    *components.SysLog
	status string
	name   string
	events map[string]func(evt *NodeServerEvent)
}

// NewNodeServer 创建一个新的主服务客户端
func NewNodeServer() *NodeServer {
	return &NodeServer{
		log:    components.NewSysLog("node_server_"),
		status: StatusClose,
		events: make(map[string]func(evt *NodeServerEvent)),
	}
}

// OnConnected 连接完成事件
func (n *NodeServer) OnConnected(e *components.TCPConnEvent) {
	n.conn = e.Conn
	n.Ping()
	n.status = StatusOpen
	if evt, ok := n.events["connected"]; ok {
		evt(&NodeServerEvent{
			Data:   e.Data,
			Client: n,
		})
	}
}

// OnDisconnected 关闭连接
func (n *NodeServer) OnDisconnected(e *components.TCPConnEvent) {
	if evt, ok := n.events["disconnect"]; ok {
		evt(&NodeServerEvent{
			Data:   e.Data,
			Client: n,
		})
	}
	n.status = StatusClose
}

// OnRecv 接收数据
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

// OnWritten 写入数据后
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

// OnError 错误事件
func (n *NodeServer) OnError(evt *components.TCPConnEvent) {
	fmt.Println("error:", evt.Data)
}

// On 外部绑定事件
func (n *NodeServer) On(evtName string, evt func(evt *NodeServerEvent)) {
	n.events[evtName] = evt
}

// Name 节点名称
func (n *NodeServer) Name() string {
	return n.name
}

// Status 节点状态
func (n *NodeServer) Status() int {
	return n.conn.Status()
}

// RemoteAddr 得到当前连接IP
func (n *NodeServer) RemoteAddr() string {
	ipstr := strings.Split(n.conn.RemoteAddr(), ":")[0]
	return ipstr
}

func (n *NodeServer) WriteData(data []byte) {
	n.conn.WriteData(data)
}

// Ping 周期性PING
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
			if evt, ok := n.events["login"]; ok {
				evt(&NodeServerEvent{
					Client: n,
				})
			}
			n.status = StatusActive
		} else {
			ackCmd := components2.NewMainStream()
			ackCmd.Command = components2.CMDClose
			n.conn.WriteData(ackCmd.Build())
		}
	case components2.CMDShell:
		shell := new(CmdShell)
		shell.Parse(cmd.Content)
		if evt, ok := n.events["ackshell"]; ok {
			evt(&NodeServerEvent{
				Client: n,
				Data:   shell,
			})
		}
	}
}
