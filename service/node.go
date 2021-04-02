package service

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"net"
	"system-monitoring/socketcon"
)

//控制节点服务
type NodeService struct {
	serverAddr string //主服务地址
	conn       *components.TCPConnect
	node       *socketcon.NodeClient
}

//初始化一个节点服务
func NewNodeService(mainAddr string) *NodeService {
	return &NodeService{
		serverAddr: mainAddr,
	}
}

//连接主服务
func (n *NodeService) Connect() error {
	conn, err := net.Dial("TCP", n.serverAddr)
	if err != nil {
		return fmt.Errorf("connect main server error: %v", err)
	}
	n.node = socketcon.NewNodeClient()
	processTcp := components.NewTCPConnect(conn, n.node)
	processTcp.Run()
	n.conn = processTcp
	return nil
}

//发送数据
func (n *NodeService) SendData(data []byte) error {
	n.conn.WriteData(data)
	return nil
}

//得到node
func (n *NodeService) Node() *socketcon.NodeClient {
	return n.node
}
