package socketcon

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"system-monitoring/common"
	components2 "system-monitoring/components"
)

//节点控制
type NodeServer struct {
	conn   *components.TCPConnect
	log    *components.SysLog
	status string
	name   string
}

//创建一个新的主服务
func NewNodeServer() *NodeServer {
	return &NodeServer{
		log: components.NewSysLog("node_client_"),
	}
}

//连接完成事件
func (n *NodeServer) OnConnected(evt *components.TCPConnEvent) {
	n.conn = evt.Conn
}

//关闭连接
func (n *NodeServer) OnDisconnected(evt *components.TCPConnEvent) {

}

//接收数据
func (n *NodeServer) OnRecv(evt *components.TCPConnEvent) {
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
}

//节点名称
func (n *NodeServer) Name() string {
	return n.name
}

//执行命令
func (n *NodeServer) execCommand(cmd *components2.MainStream) {
	switch cmd.Command {
	case components2.CMDPing:
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDAuth:
		auth := new(Auth)
		auth.Parse(cmd.Content)
		if auth.Auth == common.Conf.Server.AuthPass {
			n.name = auth.Name
		}
	}
}
