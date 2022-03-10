package socketcon

import (
	"bytes"
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/creack/pty"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"system-monitoring/common"
	components2 "system-monitoring/components"
)

// NodeClient 节点控制
type NodeClient struct {
	conn      *components.TCPConnect
	log       *components.SysLog
	status    string
	events    map[string]func(evt *components.TCPConnEvent)
	ptymx     *os.File
	ptyConn   net.Conn
	ptyIsOpen bool
}

// NewNodeClient 创建一个新的节点客户端
func NewNodeClient() *NodeClient {
	return &NodeClient{
		log:    components.NewSysLog("node_client_"),
		status: StatusOpen,
		events: make(map[string]func(evt *components.TCPConnEvent)),
	}
}

// OnConnected 连接完成事件
func (n *NodeClient) OnConnected(e *components.TCPConnEvent) {
	n.conn = e.Conn
	if evt, ok := n.events["connected"]; ok {
		e.Data = n
		evt(e)
	}
}

// OnDisconnected 关闭连接
func (n *NodeClient) OnDisconnected(e *components.TCPConnEvent) {
	if evt, ok := n.events["disconnect"]; ok {
		evt(e)
	}
}

// OnRecv 接收数据
func (n *NodeClient) OnRecv(evt *components.TCPConnEvent) {
	fmt.Println("node recv:", evt)
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
func (n *NodeClient) OnWritten(evt *components.TCPConnEvent) {
	//p.conn.Close()
}

// OnError 错误事件
func (n *NodeClient) OnError(evt *components.TCPConnEvent) {
	fmt.Println(evt.Data)
}

//执行命令
func (n *NodeClient) execCommand(cmd *components2.MainStream) {
	fmt.Println("exec command:", cmd.Command)
	switch cmd.Command {
	case components2.CMDPing: //收到服务端ping请求
		fmt.Println("server:", n.conn.RemoteAddr(), "-> ping")
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDClose: //收到服务端关闭请求
		n.conn.Close()
		n.status = StatusClose
	case components2.CMDAuthCode: //收到服务端登录成功标识
		n.status = StatusActive
		fmt.Println("login done")
		if evt, ok := n.events["login"]; ok {
			evt(&components.TCPConnEvent{
				Data: cmd,
			})
		}
	case components2.CMDPtyOpen: //收到打开pty终端
		n.openPty(cmd)
	case components2.CMDPtyClose:
		_ = n.ptymx.Close()
		_ = n.ptyConn.Close()
		n.ptymx = nil
		n.ptyConn = nil
		ptyCmd := components2.NewMainStream()
		ptyCmd.Command = components2.CMDPtyClose
		n.conn.WriteData(ptyCmd.Build())
	case components2.CMDPty:
		//fmt.Println("write pty:",string(cmd.Content))
		//_, _ = n.ptymx.Write(cmd.Content)
	}
}

// On 外部绑定事件
func (n *NodeClient) On(evtName string, evt func(evt *components.TCPConnEvent)) {
	n.events[evtName] = evt
}

//执行shell命令并返回结果
func (n *NodeClient) execShell(cmd *CmdShell) {
	shell := exec.Command(cmd.Cmd, cmd.Args...)
	shell.Dir = cmd.Dir
	var buf bytes.Buffer
	shell.Stdout = &buf
	shell.Stderr = &buf
	err := shell.Run()
	if err != nil {
		buf.WriteString(fmt.Sprintf("[ERROR] %s", err))
	}

	shellData := &CmdShell{
		AckId:      cmd.AckId,
		AckContent: buf.Bytes(),
	}
	fmt.Println("exec shell:", string(shellData.AckContent))
	ackCmd := components2.NewMainStream()
	ackCmd.Gzip(true)
	ackCmd.Command = components2.CMDShell
	ackCmd.Content = shellData.Build()
	n.conn.WriteData(ackCmd.Build())
}

//处理打开pty事件
func (n *NodeClient) openPty(cmd *components2.MainStream) {
	if n.ptyIsOpen {
		ptyCmd := components2.NewMainStream()
		ptyCmd.Command = components2.CMDPtyOpen
		n.conn.WriteData(ptyCmd.Build())
	}
	ipStr := strings.Split(n.conn.RemoteAddr(), ":")
	port, err := strconv.Atoi(ipStr[1])
	if err != nil {
		common.DebugF("parse ip port error:%v", err)
		n.sendPtyError(fmt.Errorf("parse ip port error:%v", err))
		return
	}
	ip := fmt.Sprintf("%s:%d", ipStr[0], port+1)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		common.DebugF("pty tcp connect error:%v", err)
		n.sendPtyError(fmt.Errorf("pty tcp connect error:%v", err))
		return
	}
	n.ptyConn = conn
	n.ptymx, err = GetPty()
	if err != nil {
		n.sendPtyError(fmt.Errorf("open pty error: %v", err))
		return
	}
	go n.runPty()
	n.ptyIsOpen = true
	ptyCmd := components2.NewMainStream()
	ptyCmd.Command = components2.CMDPtyOpen
	n.conn.WriteData(ptyCmd.Build())
}

func (n *NodeClient) runPty() {
	go func() {
		_, _ = io.Copy(n.ptyConn, n.ptymx)
	}()
	_, _ = io.Copy(n.ptymx, n.ptyConn)
	common.DebugF("closed pty")
}

func (n *NodeClient) sendPtyError(err error) {
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDPtyErr
	cmd.Content = []byte(err.Error())
	n.conn.WriteData(cmd.Build())
}

func GetPty() (*os.File, error) {
	// Create arbitrary command.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	return ptmx, nil
}
