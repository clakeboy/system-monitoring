package service

import (
	"encoding/json"
	"github.com/clakeboy/golib/components"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	net2 "github.com/shirou/gopsutil/net"
	"log"
	"net"
	"system-monitoring/command"
	components2 "system-monitoring/components"
	"system-monitoring/models"
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
			//DebugF("ticker time: %d", time.Now().Unix())
			data := new(models.NodeInfoData)
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Println("memory info error:", err)
				break
			}
			data.Memory = memInfo
			//fmt.Print("memory:")
			//utils.PrintAny(memInfo)

			netInterface, err := net2.Interfaces()
			//fmt.Print("net interface:")
			//utils.PrintAny(netInterface)
			data.NetInterface = netInterface
			netInfo, err := net2.IOCounters(false)
			//netInfo,err := net2.Connections("")
			if err != nil {
				log.Println("network info error:", err)
			}
			data.NetIo = netInfo
			//fmt.Print("net io:")
			//utils.PrintAny(netInfo)
			//netStats , err := net2.ConntrackStats(false)
			//fmt.Print("net stats:")
			//utils.PrintAny(netStats)
			//fmt.Println(cpu.Info())
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
			cpuInfo := &models.CpuUse{
				List: cpuList,
				All:  allCpu,
			}
			data.CpuUse = cpuInfo
			//fmt.Print("cpu use:")
			//utils.PrintAny(cpuInfo)
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
			data.DiskUse = useList
			//fmt.Print("disk use:")
			//utils.PrintAny(useList)
			diskInfo, err := disk.IOCounters()
			//diskInfo,err := disk.Usage("/")
			if err != nil {
				log.Println("disk info error:", err)
			}
			//fmt.Print("disk io:")
			//utils.PrintAny(diskInfo)
			data.DiskIo = diskInfo
			zipData, err := json.Marshal(data)
			if err != nil {
				log.Println("json info error:", err)
			}
			//zipData := bytes.NewBuffer([]byte{})
			//tw := tar.NewWriter(zipData)
			//var tmp []byte
			//tmp, _ = json.Marshal(memInfo)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "memory",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//tmp, _ = json.Marshal(netInfo)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "net_io",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//tmp, _ = json.Marshal(netInterface)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "net_interface",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//tmp, _ = json.Marshal(cpuInfo)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "cpu_use",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//tmp, _ = json.Marshal(diskInfo)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "disk_io",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//tmp, _ = json.Marshal(useList)
			//_ = tw.WriteHeader(&tar.Header{
			//	Name: "disk_use",
			//	Size: int64(len(tmp)),
			//})
			//_, _ = tw.Write(tmp)
			//_ = tw.Close()
			gData, err := components2.Gzip(zipData)
			if err != nil {
				log.Println("gzip data error:", err)
			}

			mainStream := components2.NewMainStream()
			mainStream.Command = components2.CMDSysInfo
			mainStream.Content = gData
			n.SendData(mainStream.Build())
			//
			//DebugF("org size: %d", len(zipData))
			//DebugF("gzip size: %d", len(gData))
			//unData, err := components2.UnGzip(data)
			//if err != nil {
			//	log.Println("un gzip data error:", err)
			//}
			//buf := bytes.NewReader(unData)
			//tr := tar.NewReader(buf)
			//for {
			//	th,err := tr.Next()
			//	if err == io.EOF {
			//		break
			//	}
			//	if err != nil {
			//		//return nil, err
			//		fmt.Println("read header error:",err)
			//		return
			//	}
			//	utils.PrintAny(th)
			//	content := bytes.NewBuffer([]byte{})
			//	wl,err := io.Copy(content,tr)
			//	if err != nil {
			//		//return nil, err
			//		fmt.Println("read content error:",err)
			//		return
			//	}
			//	fmt.Println(th.Name,wl)
			//	fmt.Println(content.String())
			//}
		}
	}
}

// OnDisconnect 连接断开事件
func (n *NodeService) OnDisconnect(evt *components.TCPConnEvent) {
	DebugF("disconnected for remote server: %v", evt.Conn.RemoteAddr())
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
