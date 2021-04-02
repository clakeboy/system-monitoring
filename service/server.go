package service

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"net"
	"system-monitoring/socketcon"
)

//已连接客户端列表

type TcpServer struct {
	ip    string
	debug bool
	list  []*socketcon.NodeServer
}

func NewTcpServer(ip string, debug bool) *TcpServer {
	return &TcpServer{
		ip:    ip,
		debug: debug,
	}
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
		processTcp := components.NewTCPConnect(conn, client)
		processTcp.Run()
		t.list = append(t.list, client)
	}
}

//显示连接数
func (t *TcpServer) Connected() {

}
