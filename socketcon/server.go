package socketcon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"system-monitoring/common"
	components2 "system-monitoring/components"
	"system-monitoring/models"
	"time"
)

// 事件列表
type NodeServerEventList map[string]func(data *NodeServerEvent)

// NodeServerEvent 客户端事件
type NodeServerEvent struct {
	Client *NodeServer //客户端
	Data   interface{} //事件数据
}

// NodeServer 节点控制
type NodeServer struct {
	conn      *components.TCPConnect
	log       *components.SysLog
	status    string
	name      string
	isOpenPty bool
	events    map[string]NodeServerEventList
	lock      sync.Mutex
}

// NewNodeServer 创建一个新的主服务客户端
func NewNodeServer() *NodeServer {
	return &NodeServer{
		log:    components.NewSysLog("node_server_"),
		status: StatusClose,
		events: make(map[string]NodeServerEventList),
	}
}

func (n *NodeServer) triggerEvent(evtName string, data *NodeServerEvent) {
	list, ok := n.events[evtName]
	if !ok {
		return
	}
	for _, v := range list {
		v(data)
	}
}

func (n *NodeServer) triggerEventForSn(evtName string, sn string, data *NodeServerEvent) {
	list, ok := n.events[evtName]
	if !ok {
		return
	}
	evt, ok := list[sn]
	if ok {
		evt(data)
	}
}

// On 外部绑定事件
func (n *NodeServer) On(evtName, key string, evt func(data *NodeServerEvent)) {
	n.lock.Lock()
	list, ok := n.events[evtName]
	if !ok {
		list = make(NodeServerEventList)
		n.events[evtName] = list
	}
	list[key] = evt
	n.lock.Unlock()
}

func (n *NodeServer) Off(evtName, key string, evt func(data *NodeServerEvent)) {
	n.lock.Lock()
	list, ok := n.events[evtName]
	if !ok {
		return
	}
	for k, _ := range list {
		if key == k {
			delete(list, key)
			break
		}
	}
	n.lock.Unlock()
}

// OnConnected 连接完成事件
func (n *NodeServer) OnConnected(e *components.TCPConnEvent) {
	n.conn = e.Conn
	n.Ping()
	n.status = StatusOpen
	n.triggerEvent("connected", &NodeServerEvent{
		Data:   e.Data,
		Client: n,
	})
}

// OnDisconnected 关闭连接
func (n *NodeServer) OnDisconnected(e *components.TCPConnEvent) {
	n.triggerEvent("disconnect", &NodeServerEvent{
		Data:   e.Data,
		Client: n,
	})
	n.status = StatusClose
}

// OnRecv 接收数据
func (n *NodeServer) OnRecv(evt *components.TCPConnEvent) {
	data := evt.Data.([]byte)
	if len(data) <= 0 {
		return
	}
	list, err := components2.CheckMultiStream(data)
	if err != nil {
		n.log.Error(err)
		return
	}

	if len(list) <= 0 {
		return
	}

	for _, v := range list {
		n.execCommand(v)
	}
}

// 检查是否粘包
func (n *NodeServer) checkMultiData(data []byte) [][]byte {
	var dataList [][]byte
	buf := bytes.NewBuffer([]byte{})
	read := bytes.NewBuffer(data)
	finish := false
	for {
		n, err := read.ReadBytes(0xec)
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
				common.DebugF("ping -> %s %s", n.name, n.conn.RemoteAddr())
				cmd := components2.NewMainStream()
				cmd.Command = components2.CMDPing
				n.conn.WriteData(cmd.Build())
			}
		}
	}()
}

// 执行命令
func (n *NodeServer) execCommand(cmd *components2.MainStream) {
	common.DebugF("exec command: %d", cmd.Command)
	switch cmd.Command {
	case components2.CMDPong: //节点端返回的pong返回
		common.DebugF("pong -> %s %s", n.name, n.conn.RemoteAddr())
	case components2.CMDPing: //节点端ping请求
		ackCmd := components2.NewMainStream()
		ackCmd.Command = components2.CMDPong
		n.conn.WriteData(ackCmd.Build())
	case components2.CMDAuth: //节点端登录请求
		auth := new(Auth)
		auth.Parse(cmd.Content)
		common.DebugF("login: %s", auth.Name)
		if auth.Auth == common.Conf.Server.AuthPass {
			n.name = auth.Name
			ackCmd := components2.NewMainStream()
			ackCmd.Command = components2.CMDAuthCode
			n.conn.WriteData(ackCmd.Build())
			n.triggerEvent("login", &NodeServerEvent{
				Client: n,
			})
			n.status = StatusActive
		} else {
			ackCmd := components2.NewMainStream()
			ackCmd.Command = components2.CMDClose
			n.conn.WriteData(ackCmd.Build())
		}
	case components2.CMDShell: //收到节点端返回的shell执行结果
		shell := new(CmdShell)
		shell.Parse(cmd.Content)
		n.triggerEvent("ackshell", &NodeServerEvent{
			Client: n,
			Data:   shell,
		})
	case components2.CMDPtyOpen:
		n.isOpenPty = true
		n.triggerEvent("pty_open", &NodeServerEvent{
			Client: n,
			Data:   cmd.Content,
		})
	case components2.CMDPtyClose:
		n.isOpenPty = false
		n.triggerEvent("pty_close", &NodeServerEvent{
			Client: n,
			Data:   cmd.Content,
		})

	case components2.CMDPty: //收到节点端返回的 pty信息
		n.triggerEvent("pty", &NodeServerEvent{
			Client: n,
			Data:   cmd.Content,
		})
	case components2.CMDPtyErr:
		n.triggerEvent("pty_error", &NodeServerEvent{
			Client: n,
			Data:   cmd.Content,
		})

	case components2.CMDFile: //文件处理回执
		fileInfo := new(CMDFileInfo)
		fileInfo.Parse(cmd.Content)
		model := models.NewFileModel(nil)
		data, err := model.GetById(fileInfo.FileId)
		if err != nil {
			return
		}
		if fileInfo.Error != "" {
			data.PushResult = fileInfo.Error
		} else {
			data.PushResult = fileInfo.Message
		}
		_ = model.Update(data)
	case components2.CMDSysInfo: //节点发送系统信息数据
		//model := models.NewNodeInfoModel(1)
		n.ReceiveNodeInfo(cmd)
	case components2.CMDDir: //接收数据
		dir := new(CMDDir)
		err := dir.Parse(cmd.Content)
		if err != nil {
			n.log.Error(fmt.Errorf("cmddir error:%v", err))
			return
		}
		n.triggerEventForSn(dir.Type.String(), dir.Sn, &NodeServerEvent{
			Client: n,
			Data:   dir,
		})
	}
}

// OpenPty 打开远端pty
func (n *NodeServer) OpenPty() {
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDPtyOpen
	n.WriteData(cmd.Build())
}

// ClosePty 关闭远端pty
func (n *NodeServer) ClosePty(id string) {
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDPtyClose
	cmd.Content = []byte(id)
	n.WriteData(cmd.Build())
}

// PtyIsOpen pty status 状态
func (n *NodeServer) PtyIsOpen() bool {
	return n.isOpenPty
}

// PushFile 推送数据到远程服务
func (n *NodeServer) PushFile(data *models.FileData, serv *models.ServiceData) {
	fileData := new(CMDFileInfo)
	fileData.Path = serv.Directory
	fileData.FileUri = fmt.Sprintf(":%s/serv/file/get?fid=%d", common.Conf.System.Port, data.Id)
	fileData.FileId = data.Id
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDFile
	cmd.Content = fileData.Build()
	n.WriteData(cmd.Build())
}

// ReceiveNodeInfo 收集节点服务器信息
func (n *NodeServer) ReceiveNodeInfo(cmd *components2.MainStream) {
	unData, err := components2.UnGzip(cmd.Content)
	if err != nil {
		common.DebugF("unzip node info error: %v", err)
		return
	}
	data := new(models.NodeInfoData)
	err = json.Unmarshal(unData, data)
	if err != nil {
		common.DebugF("un json data error: %v", err)
	}
	data.CreatedDate = time.Now().Unix()
	model := models.NewNodeInfoModel(n.RemoteAddr())
	if model == nil {
		return
	}

	err = model.SaveRange(data)
	if err != nil {
		common.DebugF("save data error: %v", err)
	}
}

//******************************************************
//* 文件操作相关功能

// GetRemoteDir 得到目录文件列表
func (n *NodeServer) GetRemoteDir(path string, page, number int) ([]*CMDDirList, int, error) {
	sn := utils.CreateUUID(false)
	dir := &CMDDir{
		Type:   DirList,
		Path:   path,
		Page:   page,
		Number: number,
		Sn:     sn,
	}
	wait := make(chan bool, 1)
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDDir
	cmd.Content = dir.Build()

	var list []*CMDDirList
	var count int
	var err error
	n.On(DirList.String(), sn, func(data *NodeServerEvent) {
		dirData, ok := data.Data.(*CMDDir)
		if !ok {
			wait <- true
			return
		}
		//utils.PrintAny(dirData)
		if dirData.Error != "" {
			err = fmt.Errorf("remote error: %s", dirData.Error)
			wait <- true
			return
		}
		list = dirData.List
		count = dirData.Count
		wait <- true
	})
	n.WriteData(cmd.Build())
	ctx, cancelWait := context.WithCancel(context.Background())
	go timeout(30, wait, ctx)
	<-wait
	cancelWait()
	close(wait)
	wait = nil
	//tk.Stop()
	n.Off(DirList.String(), sn, nil)
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (n *NodeServer) GetRemoteFile(path string) ([]byte, error) {
	sn := utils.CreateUUID(false)
	dir := &CMDDir{
		Path: path,
		Type: DirContent,
		Sn:   sn,
	}
	wait := make(chan bool, 1)
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDDir
	cmd.Content = dir.Build()

	var content []byte
	var err error
	n.On(DirContent.String(), sn, func(data *NodeServerEvent) {
		dirData, ok := data.Data.(*CMDDir)
		if !ok {
			wait <- true
			return
		}
		if dirData.Error != "" {
			err = fmt.Errorf("remote error: %s", dirData.Error)
			wait <- true
			return
		}
		content = dirData.Content
		wait <- true
	})
	n.WriteData(cmd.Build())
	ctx, cancelWait := context.WithCancel(context.Background())

	go timeout(30, wait, ctx)
	<-wait
	cancelWait()
	close(wait)
	wait = nil
	n.Off(DirContent.String(), sn, nil)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (n *NodeServer) SaveRemoteFile(path string, con []byte) error {
	sn := utils.CreateUUID(false)
	dir := &CMDDir{
		Path:    path,
		Type:    DirSaveFile,
		Content: con,
		Sn:      sn,
	}
	wait := make(chan bool, 1)
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDDir
	cmd.Content = dir.Build()

	var err error
	n.On(DirSaveFile.String(), sn, func(data *NodeServerEvent) {
		dirData, ok := data.Data.(*CMDDir)
		if !ok {
			wait <- true
			return
		}

		if dirData.Error != "" {
			err = fmt.Errorf("remote error: %s", dirData.Error)
			wait <- true
			return
		}
		wait <- true
	})
	n.WriteData(cmd.Build())
	ctx, cancelWait := context.WithCancel(context.Background())

	go timeout(30, wait, ctx)
	<-wait
	cancelWait()
	close(wait)
	wait = nil
	n.Off(DirSaveFile.String(), sn, nil)
	if err != nil {
		return err
	}
	return nil
}

func timeout(sec int, wait chan bool, ctx context.Context) {
	tk := time.NewTicker(time.Second * time.Duration(sec))
	loop := true
	for loop {
		select {
		case <-tk.C:
			loop = false
		case <-ctx.Done():
			tk.Stop()
			return
		}
	}
	tk.Stop()
	wait <- true
}
