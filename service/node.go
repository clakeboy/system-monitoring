package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"system-monitoring/command"
	components2 "system-monitoring/components"
	"system-monitoring/models"
	"system-monitoring/socketcon"
	"time"

	"github.com/clakeboy/golib/components"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	net2 "github.com/shirou/gopsutil/net"
)

// NodeService 控制节点服务
type NodeService struct {
	serverAddr      string //主服务地址
	conn            *components.TCPConnect
	node            *socketcon.NodeClient
	mode            string    //连接模式
	reconnectNumber int       //重新连接次数
	name            string    //节点名称
	passwd          string    //验证密码
	stopInfo        chan bool //发送停止信号
}

// NewNodeService 初始化一个节点服务
func NewNodeService(mainAddr, name, passwd string) *NodeService {
	return &NodeService{
		serverAddr: mainAddr,
		name:       name,
		passwd:     passwd,
		stopInfo:   make(chan bool, 1),
	}
}

// PassiveConnect 开启被动连接
func (n *NodeService) PassiveConnect() error {
	n.mode = "passive"
	server, err := net.Listen("tcp", ":17511")
	if err != nil {
		return err
	}
	defer server.Close()
	conn, err := server.Accept()
	if err != nil {
		return err
	}
	n.initService(conn)
	return nil
}

// Connect 主动连接主服务
func (n *NodeService) Connect() error {
	n.mode = "active"
	conn, err := net.Dial("tcp", n.serverAddr)
	if err != nil {
		DebugF("connect main server error: %v", err)
		DebugF("10s will be reconnect to: %s", n.serverAddr)
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

//func (n *NodeService) PtyConnect() error {
//	server, err := net.Listen("tcp", ":17511")
//	if err != nil {
//		return err
//	}
//	defer server.Close()
//	conn, err := server.Accept()
//	if err != nil {
//		return err
//	}
//	n.ptyConn = conn
//	return nil
//}

func (n *NodeService) initService(conn net.Conn) {
	n.node = socketcon.NewNodeClient()
	n.node.On("disconnect", n.OnDisconnect)
	n.node.On("login", n.OnLogin)
	//n.node.On("pty", n.OnPty)
	processTcp := components.NewTCPConnect(conn, n.node)
	processTcp.Run()
	processTcp.SetReadTimeout(time.Minute * 5)
	//processTcp.SetWriteTimeout(0)
	processTcp.SetDebug(command.CmdDebug)
	n.conn = processTcp
	//发送登录验证信息
	n.Auth()
}

// Auth 登录系统
func (n *NodeService) Auth() {
	auth := new(socketcon.Auth)
	auth.Name = n.name
	auth.Auth = n.passwd
	cmd := components2.NewMainStream()
	cmd.Command = components2.CMDAuth
	cmd.Content = auth.Build()
	n.SendData(cmd.Build())
}

// SendData 发送数据
func (n *NodeService) SendData(data []byte) {
	n.conn.WriteData(data)
}

// Node 得到node
func (n *NodeService) Node() *socketcon.NodeClient {
	return n.node
}

func (n *NodeService) sendServerStatus() {
	tk := time.NewTicker(time.Second * 5)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			data := new(models.NodeInfoData)
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Println("memory info error:", err)
				break
			}
			data.Memory = memInfo

			netInterface, err := net2.Interfaces()
			data.NetInterface = netInterface
			netInfo, err := net2.IOCounters(false)
			if err != nil {
				log.Println("network info error:", err)
			}
			data.NetIo = netInfo
			cpuList, err := cpu.Percent(0, true)
			if err != nil {
				log.Println("cpu info error:", err)
			}
			allCpu := 0.0
			for _, v := range cpuList {
				allCpu += v
			}
			allCpu = allCpu / float64(len(cpuList))
			cpuInfo := &models.CpuUse{
				List: cpuList,
				All:  allCpu,
			}
			data.CpuUse = cpuInfo
			part, err := disk.Partitions(false)
			if err != nil {
				log.Println("disk part error:", err)
			}
			var useList []*disk.UsageStat
			for _, v := range part {
				usage, err := disk.Usage(v.Mountpoint)
				if err != nil {
					continue
				}
				useList = append(useList, usage)
			}
			data.DiskUse = useList
			diskInfo, err := disk.IOCounters()
			if err != nil {
				log.Println("disk info error:", err)
			}
			data.DiskIo = diskInfo
			zipData, err := json.Marshal(data)
			if err != nil {
				log.Println("json info error:", err)
			}
			gData, err := components2.Gzip(zipData)
			if err != nil {
				log.Println("gzip data error:", err)
			}

			mainStream := components2.NewMainStream()
			mainStream.Command = components2.CMDSysInfo
			mainStream.Content = gData

			n.SendData(mainStream.Build())
		case <-n.stopInfo:
			fmt.Printf("stop send info:%d\n", time.Now().Unix())
			return
		}
	}
}

// OnDisconnect 连接断开事件
func (n *NodeService) OnDisconnect(evt *components.TCPConnEvent) {
	DebugF("disconnected for remote server: %v", evt.Conn.RemoteAddr())
	n.stopInfo <- true
	if n.mode == "passive" {
		_ = n.PassiveConnect()
	} else {
		_ = n.Connect()
	}
}

// OnLogin 登录成功后开始发送状态信息
func (n *NodeService) OnLogin(evt *components.TCPConnEvent) {
	if !command.CmdDebug {
		go n.sendServerStatus()
	}
}

//func (n *NodeService) OnPty(evt *components.TCPConnEvent) {
//	err := n.PtyConnect()
//	if err != nil {
//		return
//	}
//	pymx,err := socketcon.GetPty()
//	if err != nil {
//		return
//	}
//	go io.Copy(n.ptyConn,pymx)
//	_,_ = io.Copy(pymx,n.ptyConn)
//}
