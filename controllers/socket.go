package controllers

import (
	"fmt"
	"io"
	"net"
	"system-monitoring/common"
	components2 "system-monitoring/components"
	"system-monitoring/models"
	"system-monitoring/service"
	"system-monitoring/socketcon"
	"system-monitoring/websocket"
)

type SocketController struct {
	so          *websocket.Client
	cancelClone bool                  //是否取消克隆
	nodeServer  *socketcon.NodeServer //节点服务
	ptyConn     net.Conn              //远程pty连接
	key         string
}

func NewSocketController() *SocketController {
	return &SocketController{}
}

//绑定socket 事件处理器
func (s *SocketController) Connect(so *websocket.Client) error {
	s.so = so
	s.key = so.Id()
	so.On(models.SocketEventPty, s.onPty)
	so.On("close", s.onClose)
	return nil
}

func (s *SocketController) onPty(data []byte) []byte {
	ptyData := models.NewPtyMessage()
	err := ptyData.ParseJson(data)
	if err != nil {
		return []byte(fmt.Sprintf("parse data error %s", err))
	}

	switch ptyData.Evt {
	case models.PtyStart:
		fmt.Println("controler start", s.key)
		if s.nodeServer != nil && s.nodeServer.PtyIsOpen() {
			return nil
		}
		data, err := models.ParseNodeData(ptyData.Data)
		if err != nil {
			return []byte(fmt.Sprintf("wrong node data %s", err))
		}
		s.nodeServer, err = service.MainServer.GetNodeServer(data.Ip)
		if err != nil {
			return []byte(fmt.Sprintf("node server error %s", err))
		}
		s.nodeServer.On("pty_error", s.key, s.ptyError)
		s.nodeServer.On("pty_close", s.key, s.ptyClose)
		s.nodeServer.On("pty_open", s.key, s.ptyOpen)
		s.nodeServer.OpenPty()
		return []byte("opening remote pty...")
	case models.PtyExec:
		if s.nodeServer == nil || s.ptyConn == nil {
			return []byte("pls open remote pty...\r\n")
		}
		_, err := s.ptyConn.Write([]byte(ptyData.Cmd))
		if err != nil {
			return []byte(fmt.Sprintf("send pty data error: %v", err))
		}
	case models.PtyClose:
		if s.nodeServer == nil || s.ptyConn == nil {
			return []byte("pls open remote pty...\r\n")
		}
		s.nodeServer.ClosePty(s.key)
	}
	return nil
}

func (s *SocketController) ptyOpen(evt *socketcon.NodeServerEvent) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("pty open error: ", err)
		}
	}()
	portStr := string(evt.Data.([]byte))
	ipStr := s.nodeServer.RemoteAddr()
	//conn, err := service.MainServer.GetPtyConn(s.nodeServer.RemoteAddr())
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ipStr, portStr))
	if err != nil {
		_ = s.so.Emit(models.SocketEventPty, []byte(err.Error()), nil)
		return
	}

	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDPty
	cmd.Content = []byte(s.key)
	_, err = conn.Write(cmd.Build())
	if err != nil {
		_ = s.so.Emit(models.SocketEventPty, []byte(err.Error()), nil)
		return
	}
	s.ptyConn = conn
	go func() {
		_, _ = io.Copy(s, s.ptyConn)
		fmt.Println("close terminal")
	}()
}

func (s *SocketController) Write(data []byte) (int, error) {
	if !s.nodeServer.PtyIsOpen() {
		_ = s.so.Emit(models.SocketEventPty, []byte("remote pty is not open"), nil)
		return 0, fmt.Errorf("remote pty is not open")
	}
	err := s.so.Emit(models.SocketEventPty, data, nil)
	if err != nil {
		fmt.Println("push pty data error", err)
		return 0, err
	}
	return len(data), nil
}

func (s *SocketController) ptyError(evt *socketcon.NodeServerEvent) {
	common.DebugF("receive pty error data:%s", evt.Data)
	err := s.so.Emit(models.SocketEventPty, append(evt.Data.([]byte), '\r', '\n'), nil)
	if err != nil {
		fmt.Println("push pty data error", err)
	}
}

func (s *SocketController) ptyClose(evt *socketcon.NodeServerEvent) {
	if s.ptyConn != nil {
		_ = s.ptyConn.Close()
	}
	s.ptyConn = nil
	s.nodeServer = nil
	err := s.so.Emit(models.SocketEventPty, []byte("remote pty is closed!\n"), nil)
	if err != nil {
		fmt.Println("push pty data error", err)
	}
}

func (s *SocketController) onClose() {
	fmt.Println("close web socket")
	if s.nodeServer != nil {
		s.nodeServer.Off("pty_error", s.key, nil)
		s.nodeServer.Off("pty_close", s.key, nil)
		s.nodeServer.Off("pty_open", s.key, nil)
		s.nodeServer.ClosePty(s.key)
	}
	s.ptyClose(nil)
}

////导入事件
//func (s *SocketController) onImport(data []byte) []byte {
//	params := models.NewImportParams()
//	err := params.ParseJson(data)
//	if err != nil {
//		return models.NewSocketResult(models.ErrorCode, err.Error(), nil).ToJson()
//	}
//	if params.Code == models.CancelCode {
//		//s.cancelClone = true
//		models.CancelImport = true
//		return models.NewSocketResult(models.CancelCode, "import canceled", nil).ToJson()
//	}
//	conn, err := common.Conns.Get(params.Server.ServerId)
//	if err != nil {
//		return models.NewSocketResult(models.ErrorCode, err.Error(), nil).ToJson()
//	}
//	dbImport := models.NewImport(params, conn, s.so)
//	return dbImport.Import()
//}
