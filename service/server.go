package service

import (
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"net"
	"strconv"
	"strings"
	"sync"
	"system-monitoring/command"
	components2 "system-monitoring/components"
	"system-monitoring/models"
	"system-monitoring/socketcon"
	"time"
)

//已连接客户端列表

type TcpServer struct {
	ip    string
	debug bool
	list  map[string]*socketcon.NodeServer
	evlk  sync.RWMutex
}

func NewTcpServer(ip string, debug bool) *TcpServer {
	return &TcpServer{
		ip:    ip,
		debug: debug,
		list:  make(map[string]*socketcon.NodeServer),
	}
}

// Connect 发起主动连接节点
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
		client.On("disconnect", t.evtDisconnect)
		client.On("login", t.evtLogin)
		client.On("ackshell", t.evtAckShell)
		processTcp := components.NewTCPConnect(conn, client)
		processTcp.Run()
		processTcp.SetDebug(t.debug)
		processTcp.SetReadTimeout(time.Minute * 5)
		//processTcp.SetWriteTimeout(0)
		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
		t.evlk.Lock()
		t.list[ip] = client
		t.evlk.Unlock()
	}
}

// Connected 显示连接数
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

// CheckIp 检查ip的服务节点是否存在
func (t *TcpServer) CheckIp(ip string) bool {
	_, ok := t.list[ip]
	return ok
}

//连接断开事件
func (t *TcpServer) evtDisconnect(evt *socketcon.NodeServerEvent) {
	ipStr := evt.Client.RemoteAddr()
	fmt.Println("disconnected for node server:", ipStr)
	t.evlk.Lock()
	delete(t.list, evt.Client.RemoteAddr())
	t.evlk.Unlock()
	model := models.NewNodeModel(nil)
	data, err := model.GetByIp(ipStr)
	if err != nil {
		return
	}
	data.Status = models.NodeStatusOffline
	_ = model.Update(data)
}

//连接事件
func (t *TcpServer) evtLogin(evt *socketcon.NodeServerEvent) {
	client := evt.Client
	model := models.NewNodeModel(nil)
	data, err := model.GetByIp(client.RemoteAddr())
	if err != nil {
		data = &models.NodeData{
			Id:             0,
			Name:           client.Name(),
			Ip:             client.RemoteAddr(),
			Status:         models.NodeStatusOnline,
			LastOnlineDate: time.Now().Unix(),
			CreateDate:     time.Now().Unix(),
		}
		err = model.Save(data)
		if err != nil {
			DebugF(err.Error())
		}
	} else {
		data.LastOnlineDate = time.Now().Unix()
		data.Status = models.NodeStatusOnline
		data.Ip = client.RemoteAddr()
		err = model.Update(data)
		if err != nil {
			DebugF(err.Error())
		}
	}
}

// ExecShell 执行shell命令
func (t *TcpServer) ExecShell(ip string, cmdData *models.ShellData) error {
	node, ok := t.list[ip]
	if !ok {
		return fmt.Errorf("node server offline ")
	}
	shell := new(socketcon.CmdShell)
	shell.Cmd = cmdData.Cmd
	shell.Args = cmdData.Args
	shell.Dir = cmdData.Dir
	shell.AckId = fmt.Sprintf("%d", cmdData.Id)
	sendCmd := components2.NewMainStream()
	sendCmd.Command = components2.CMDShell
	sendCmd.Content = shell.Build()

	node.WriteData(sendCmd.Build())
	return nil
}

func (t *TcpServer) evtAckShell(evt *socketcon.NodeServerEvent) {
	data := evt.Data.(*socketcon.CmdShell)
	ackId, err := strconv.Atoi(data.AckId)
	if err != nil {
		DebugF(err.Error())
		return
	}
	model := models.NewShellModel(nil)
	shellData, err := model.GetById(ackId)
	if err != nil {
		DebugF(err.Error())
		return
	}

	shellData.ExecContent = string(data.AckContent)
	shellData.Status = 1
	err = model.Save(shellData)
	if err != nil {
		DebugF(err.Error())
		return
	}
}

var debugLog = components.NewSysLog("tcp_server_debug_")

func DebugF(str string, args ...interface{}) {
	if command.CmdDebug {
		fmt.Printf(str+"\n", args...)
	} else {
		debugLog.Info("[DEBUG] " + fmt.Sprintf(str+"\n", args...))
	}
}
