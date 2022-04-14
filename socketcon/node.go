package socketcon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"github.com/creack/pty"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
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
	ptyIsOpen bool
	ptyList   map[string]net.Conn
	ptyServer net.Listener
}

// NewNodeClient 创建一个新的节点客户端
func NewNodeClient() *NodeClient {
	return &NodeClient{
		log:     components.NewSysLog("node_client_"),
		status:  StatusOpen,
		events:  make(map[string]func(evt *components.TCPConnEvent)),
		ptyList: make(map[string]net.Conn),
	}
}

// OnConnected 连接完成事件
func (n *NodeClient) OnConnected(e *components.TCPConnEvent) {
	n.conn = e.Conn
	if evt, ok := n.events["connected"]; ok {
		e.Data = n
		evt(e)
	}
	go func() {
		_ = n.ListenPty()
	}()
}

// OnDisconnected 关闭连接
func (n *NodeClient) OnDisconnected(e *components.TCPConnEvent) {
	if evt, ok := n.events["disconnect"]; ok {
		evt(e)
	}
}

// OnRecv 接收数据
func (n *NodeClient) OnRecv(evt *components.TCPConnEvent) {
	data := evt.Data.([]byte)
	if len(data) <= 0 {
		return
	}
	list := n.checkMultiData(data)

	if len(list) <= 0 {
		return
	}

	for _, v := range list {
		common.DebugF("receive: %x", v)
		cmd := components2.NewMainStream()
		err := cmd.Parse(v)
		if err != nil {
			n.log.Error(fmt.Errorf("note receive data error: %v,\n %x", err, v))
			return
		}
		n.execCommand(cmd)
	}
}

//检查是否粘包
func (n *NodeClient) checkMultiData(data []byte) [][]byte {
	var dataList [][]byte
	buf := bytes.NewBuffer([]byte{})
	read := bytes.NewBuffer(data)
	finish := false
	for {
		n, err := read.ReadBytes(0xEC)
		if err != nil {
			break
		}

		buf.Write(n)
		if finish && bytes.Equal(n[len(n)-2:], components2.Mask) {
			dataList = append(dataList, buf.Bytes())
			buf = bytes.NewBuffer([]byte{})
			finish = false
			continue
		}

		if len(n) == 2 && !finish {
			finish = true
		}
	}
	return dataList
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
	common.DebugF("exec command: %d", cmd.Command)
	switch cmd.Command {
	case components2.CMDPing: //收到服务端ping请求
		common.DebugF("server: %s %s", n.conn.RemoteAddr(), "-> ping")
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDClose: //收到服务端关闭请求
		n.conn.Close()
		n.status = StatusClose
	case components2.CMDAuthCode: //收到服务端登录成功标识
		n.status = StatusActive
		common.DebugF("login done")
		if evt, ok := n.events["login"]; ok {
			evt(&components.TCPConnEvent{
				Data: cmd,
			})
		}
	case components2.CMDShell:
		shell := new(CmdShell)
		shell.Parse(cmd.Content)
		n.execShell(shell)
	case components2.CMDPtyOpen: //收到打开pty终端
		n.openPty(cmd)
	case components2.CMDPtyClose:
		n.closePty(cmd)
	case components2.CMDPty:
		//fmt.Println("write pty:",string(cmd.Content))
		//_, _ = n.ptymx.Write(cmd.Content)
	case components2.CMDFile:
		fileInfo := new(CMDFileInfo)
		fileInfo.Parse(cmd.Content)
		n.downloadFile(fileInfo)
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
	ptyCmd := components2.NewMainStream()
	ptyCmd.Command = components2.CMDPtyOpen
	ptyCmd.Content = []byte(common.Conf.Node.PtyPort)
	n.conn.WriteData(ptyCmd.Build())
}

//func (n *NodeClient) openPty(cmd *components2.MainStream) {
//	if n.ptyIsOpen {
//		ptyCmd := components2.NewMainStream()
//		ptyCmd.Command = components2.CMDPtyOpen
//		ptyCmd.Content = []byte(common.Conf.Node.PtyPort)
//		n.conn.WriteData(ptyCmd.Build())
//	}
//	ipStr := strings.Split(n.conn.RemoteAddr(), ":")
//	port, err := strconv.Atoi(ipStr[1])
//	if err != nil {
//		common.DebugF("parse ip port error:%v", err)
//		n.sendPtyError(fmt.Errorf("parse ip port error:%v", err))
//		return
//	}
//	ip := fmt.Sprintf("%s:%d", ipStr[0], port+1)
//	conn, err := net.Dial("tcp", ip)
//	if err != nil {
//		common.DebugF("pty tcp connect error:%v", err)
//		n.sendPtyError(fmt.Errorf("pty tcp connect error:%v", err))
//		return
//	}
//	n.ptyConn = conn
//	n.ptymx, err = GetPty()
//	if err != nil {
//		n.sendPtyError(fmt.Errorf("open pty error: %v", err))
//		return
//	}
//	go n.runPty()
//	n.ptyIsOpen = true
//	ptyCmd := components2.NewMainStream()
//	ptyCmd.Command = components2.CMDPtyOpen
//	n.conn.WriteData(ptyCmd.Build())
//}

// ListenPty 打开pty端口监听
func (n *NodeClient) ListenPty() error {
	defer func() {
		if err := recover(); err != nil {
			n.log.Error(err)
			n.ptyIsOpen = false
		}
	}()
	var err error
	n.ptyServer, err = net.Listen("tcp", fmt.Sprintf(":%s", common.Conf.Node.PtyPort))
	if err != nil {
		return err
	}
	n.ptyIsOpen = true
	for {
		client, err := n.ptyServer.Accept()
		if err != nil {
			n.ptyIsOpen = false
			return err
		}
		go n.runPty(client)
	}
}

func (n *NodeClient) runPty(client net.Conn) {
	id, err := n.checkPty(client)
	if err != nil {
		_, _ = client.Write([]byte(err.Error()))
		_ = client.Close()
		return
	}
	n.ptyList[id] = client
	ptymx, err := GetPty()
	if err != nil {
		_, _ = client.Write([]byte(err.Error()))
	}
	defer func() {
		_ = ptymx.Close()
		_ = client.Close()
	}()

	go func() {
		_, _ = io.Copy(client, ptymx)
	}()
	_, _ = io.Copy(ptymx, client)
	common.DebugF("closed pty %s", id)
	delete(n.ptyList, id)
}

func (n *NodeClient) checkPty(client net.Conn) (string, error) {
	buf := make([]byte, 1024)
	lens, err := client.Read(buf)
	if err != nil {
		return "", err
	}
	head := components2.NewMainStream()
	err = head.Parse(buf[:lens])
	if err != nil {
		return "", err
	}
	if head.Command != components2.CMDPty {
		return "", fmt.Errorf("check valid error")
	}
	return string(head.Content), nil
}

func (n *NodeClient) closePty(cmd *components2.MainStream) {
	id := string(cmd.Content)
	for k, v := range n.ptyList {
		if k == id {
			_ = v.Close()
			delete(n.ptyList, k)
			break
		}
	}
	//if n.ptymx != nil {
	//	_ = n.ptymx.Close()
	//}
	//if n.ptyConn != nil {
	//	_ = n.ptyConn.Close()
	//}
	//n.ptymx = nil
	//n.ptyConn = nil
	//n.ptyIsOpen = false
	ptyCmd := components2.NewMainStream()
	ptyCmd.Command = components2.CMDPtyClose
	n.conn.WriteData(ptyCmd.Build())
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

//******************************
//文件处理函数

//下载文件并替换
func (n *NodeClient) downloadFile(data *CMDFileInfo) {
	host := strings.Split(common.Conf.Node.Server, ":")[0]
	urlStr := fmt.Sprintf("http://%s%s", host, data.FileUri)
	client := utils.NewHttpClient()
	res, err := client.Request("GET", urlStr, nil)
	if err != nil {
		data.Error = fmt.Sprintf("http request error %v", err)
		n.pushFileResponse(data)
		return
	}

	if res.StatusCode != 200 {
		data.Error = fmt.Sprintf("request error code: %d", res.StatusCode)
		data.Message = string(res.Content)
		n.pushFileResponse(data)
		return
	}
	reg := regexp.MustCompile(`filename="(.+)"`)
	if !reg.MatchString(res.Headers.Get("Content-disposition")) {
		errMsg := utils.M{}
		err = json.Unmarshal(res.Content, &errMsg)
		if err != nil {
			data.Error = fmt.Sprintf("request error code: %d", res.StatusCode)
			data.Message = string(res.Content)
		} else {
			data.Error = fmt.Sprintf("code:%v,msg:%v", errMsg["errcode"], errMsg["errmsg"])
			data.Message = string(res.Content)
		}
		n.pushFileResponse(data)
		return
	}

	list := reg.FindStringSubmatch(res.Headers.Get("Content-disposition"))
	saveName := list[1]
	fullPath := fmt.Sprintf("%s/%s", data.Path, saveName)
	if utils.Exist(fullPath) {
		err := os.Remove(fullPath)
		if err != nil {
			data.Error = fmt.Sprintf("remove exist file error: %v", err)
			n.pushFileResponse(data)
			return
		}
	}

	err = ioutil.WriteFile(fullPath, res.Content, 0755)
	if err != nil {
		data.Error = fmt.Sprintf("write file error: %v", err)
		n.pushFileResponse(data)
		return
	}
	data.Error = ""
	data.Message = "success"
	n.pushFileResponse(data)
}

func (n *NodeClient) pushFileResponse(data *CMDFileInfo) {
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDFile
	cmd.Content = data.Build()
	n.conn.WriteData(cmd.Build())
}
