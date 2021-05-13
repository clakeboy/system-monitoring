package service

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"net"
	"system-monitoring/socketcon"
	"time"
)

//已连接客户端列表

type TcpServer struct {
	ip    string
	debug bool
	list  map[string]*socketcon.NodeServer
}

func NewTcpServer(ip string, debug bool) *TcpServer {
	return &TcpServer{
		ip:    ip,
		debug: debug,
		list:  make(map[string]*socketcon.NodeServer),
	}
}

//发起主动连接节点
func (t *TcpServer) Connect(addr string) {

}

func (t *TcpServer) Start() {
	go t.run()
}

func (t *TcpServer) run() {
	tcp, err := net.Listen("tcp", t.ip)
	if err != nil {
		panic(err)
	}

	defer tcp.Close()
	if t.debug {
		fmt.Println("Listening TCP ", t.ip)
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			components.NewSysLog("tcp_error_").Error(err)
			panic(err)
		}
		client := socketcon.NewNodeServer()
		client.On("disconnect", t.OnDisconnect)
		processTcp := components.NewTCPConnect(conn, client)
		processTcp.Run()
		processTcp.SetDebug(t.debug)
		processTcp.SetReadTimeout(time.Minute * 5)
		//processTcp.SetWriteTimeout(0)
		t.list[conn.RemoteAddr().String()] = client
	}
}

//显示连接数
func (t *TcpServer) Connected() []utils.M {
	var list []utils.M
	for k, v := range t.list {
		list = append(list, utils.M{
			"addr": k,
			"name": v.Name(),
		})
	}
	return list
}

//连接断开事件
func (t *TcpServer) OnDisconnect(evt *components.TCPConnEvent) {
	fmt.Println("disconnected for node server:", evt.Conn.RemoteAddr())
	delete(t.list, evt.Conn.RemoteAddr())
}
