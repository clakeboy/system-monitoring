package service

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"net"
	"system-monitoring/command"
	components2 "system-monitoring/components"
	"system-monitoring/socketcon"
	"time"
)

//控制节点服务
type NodeService struct {
	serverAddr      string //主服务地址
	conn            *components.TCPConnect
	node            *socketcon.NodeClient
	mode            string //连接模式
	reconnectNumber int    //重新连接次数
	name            string //节点名称
	passwd          string //验证密码
}

//初始化一个节点服务
func NewNodeService(mainAddr, name, passwd string) *NodeService {
	return &NodeService{
		serverAddr: mainAddr,
		name:       name,
		passwd:     passwd,
	}
}

//开启被动连接
func (n *NodeService) PassiveConnect() error {
	n.mode = "passive"
	server, err := net.Listen("tcp", ":17511")
	if err != nil {
		return err
	}
	conn, err := server.Accept()
	if err != nil {
		return err
	}
	n.initService(conn)
	return nil
}

//主动连接主服务
func (n *NodeService) Connect() error {
	n.mode = "active"
	conn, err := net.Dial("tcp", n.serverAddr)
	if err != nil {
		fmt.Println("connect main server error:", err)
		fmt.Println("10s will be reconnect to:", n.serverAddr)
		go func() {
			time.Sleep(time.Second * 10)
			n.reconnectNumber++
			_ = n.Connect()
		}()
		return err
	}
	n.initService(conn)
	return nil
}

func (n *NodeService) initService(conn net.Conn) {
	n.node = socketcon.NewNodeClient()
	n.node.On("disconnect", n.OnDisconnect)
	processTcp := components.NewTCPConnect(conn, n.node)
	processTcp.Run()
	processTcp.SetReadTimeout(time.Minute * 5)
	//processTcp.SetWriteTimeout(0)
	processTcp.SetDebug(command.CmdDebug)
	n.conn = processTcp
	//发送登录验证信息
	n.Auth()
}

//登录系统
func (n *NodeService) Auth() {
	auth := new(socketcon.Auth)
	auth.Name = n.name
	auth.Auth = n.passwd
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDAuth
	cmd.Content = auth.Build()
	n.SendData(cmd.Build())
}

//发送数据
func (n *NodeService) SendData(data []byte) {
	n.conn.WriteData(data)
}

//得到node
func (n *NodeService) Node() *socketcon.NodeClient {
	return n.node
}

//连接断开事件
func (n *NodeService) OnDisconnect(evt *components.TCPConnEvent) {
	fmt.Println("disconnected for remote server:", evt.Conn.RemoteAddr())
	if n.mode == "passive" {
		_ = n.PassiveConnect()
	} else {
		_ = n.Connect()
	}
}
