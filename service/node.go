package service

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	net2 "github.com/shirou/gopsutil/net"
	"log"
	"net"
	"system-monitoring/command"
	components2 "system-monitoring/components"
	"system-monitoring/socketcon"
	"time"
)

// NodeService 控制节点服务
type NodeService struct {
	serverAddr      string //主服务地址
	conn            *components.TCPConnect
	node            *socketcon.NodeClient
	mode            string //连接模式
	reconnectNumber int    //重新连接次数
	name            string //节点名称
	passwd          string //验证密码
}

// NewNodeService 初始化一个节点服务
func NewNodeService(mainAddr, name, passwd string) *NodeService {
	return &NodeService{
		serverAddr: mainAddr,
		name:       name,
		passwd:     passwd,
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
	n.node.On("login", n.OnLogin)
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
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Println("memory info error:", err)
				break
			}
			//utils.PrintAny(memInfo)
			netInfo, err := net2.IOCounters(false)
			//netInfo,err := net2.Connections("")
			if err != nil {
				log.Println("network info error:", err)
			}
			//utils.PrintAny(netInfo)
			//utils.PrintAny(netInfo)
			cpuList, err := cpu.Percent(0, true)
			if err != nil {
				log.Println("cpu info error:", err)
			}
			//utils.PrintAny(cpuList)
			allCpu := 0.0
			for _, v := range cpuList {
				allCpu += v
			}
			allCpu = allCpu / float64(len(cpuList))
			//fmt.Println(allCpu)
			cpuInfo := utils.M{
				"list": cpuList,
				"all":  allCpu,
			}
			utils.PrintAny(cpuInfo)
			part, err := disk.Partitions(false)
			if err != nil {
				log.Println("disk part error:", err)
			}
			//utils.PrintAny(part)
			var useList []*disk.UsageStat
			for _, v := range part {
				usage, err := disk.Usage(v.Mountpoint)
				if err != nil {
					continue
				}
				useList = append(useList, usage)
			}
			//utils.PrintAny(useList)
			diskInfo, err := disk.IOCounters()
			//diskInfo,err := disk.Usage("/")
			if err != nil {
				log.Println("disk info error:", err)
			}
			//utils.PrintAny(diskInfo)

			zipData := bytes.NewBuffer([]byte{})
			tw := tar.NewWriter(zipData)
			var tmp []byte
			tmp, _ = json.Marshal(memInfo)
			tw.WriteHeader(&tar.Header{
				Name: "mem",
				Size: int64(len(tmp)),
			})
			tw.Write(tmp)
			tmp, _ = json.Marshal(netInfo)
			tw.WriteHeader(&tar.Header{
				Name: "net",
				Size: int64(len(tmp)),
			})
			tw.Write(tmp)
			tmp, _ = json.Marshal(cpuInfo)
			tw.WriteHeader(&tar.Header{
				Name: "cpu",
				Size: int64(len(tmp)),
			})
			tw.Write(tmp)
			tmp, _ = json.Marshal(diskInfo)
			tw.WriteHeader(&tar.Header{
				Name: "disk",
				Size: int64(len(tmp)),
			})
			tw.Write(tmp)
			tw.Close()
			data, err := components2.Gzip(zipData.Bytes())
			if err != nil {
				log.Println("gzip data error:", err)
			}
			fmt.Println("org size:", zipData.Len())
			fmt.Println("gzip size:", len(data))
			//unData, err := components2.UnGzip(data)
			//if err != nil {
			//	log.Println("un gzip data error:", err)
			//}
		}
	}
}

// OnDisconnect 连接断开事件
func (n *NodeService) OnDisconnect(evt *components.TCPConnEvent) {
	fmt.Println("disconnected for remote server:", evt.Conn.RemoteAddr())
	if n.mode == "passive" {
		_ = n.PassiveConnect()
	} else {
		_ = n.Connect()
	}
}

// OnLogin 登录成功后开始发送状态信息
func (n *NodeService) OnLogin(evt *components.TCPConnEvent) {
	//go n.sendServerStatus()
}
