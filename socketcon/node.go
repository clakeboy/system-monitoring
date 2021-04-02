package socketcon

import "github.com/clakeboy/golib/components"

//节点控制
type NodeClient struct {
	conn   *components.TCPConnect
	log    *components.SysLog
	status string
}

//创建一个新的节点客户端
func NewNodeClient() *NodeClient {
	return &NodeClient{
		log: components.NewSysLog("node_client_"),
	}
}

//连接完成事件
func (n *NodeClient) OnConnected(evt *components.TCPConnEvent) {
	n.conn = evt.Conn
}

//关闭连接
func (n *NodeClient) OnDisconnected(evt *components.TCPConnEvent) {

}

//接收数据
func (n *NodeClient) OnRecv(evt *components.TCPConnEvent) {

}

//写入数据后
func (n *NodeClient) OnWritten(evt *components.TCPConnEvent) {
	//p.conn.Close()
}
