package main

//./sys-monitor_darwin --debug --server --config=../dev.conf
import (
	"embed"
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"os"
	"path"
	"system-monitoring/command"
	"system-monitoring/common"
	"system-monitoring/router"
	"system-monitoring/service"
	"system-monitoring/websocket"
)

var sigs chan os.Signal
var done chan bool
var (
	AppName      string //应用名称
	AppVersion   string //应用名称
	BuildVersion string //编译版本
	BuildTime    string //编译时间
	GitRevision  string //Git 版本
	GitBranch    string //Git 分支
	GoVersion    string //Golang 信息
)

//go:embed assets/html/*
var htmlFiles embed.FS

var httpServer *router.HttpServer

//var MainServer *service.TcpServer
//var NodeServer *service.NodeService

func main() {
	initService()
	go utils.ExitApp(sigs, func(s os.Signal) {
		_ = os.Remove(command.CmdPidName)
		done <- true
	})
	if command.CmdServer {
		service.MainServer.Start()
		httpServer.Start()
	} else {
		var err error
		if command.CmdPassive {
			err = service.NodeServer.PassiveConnect()
		} else {
			err = service.NodeServer.Connect()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		<-done
	}
}

func initService() {
	var err error
	command.InitCommand()
	if command.CmdShowVersion {
		Version()
		return
	}
	//获取YAML
	common.Conf = common.NewYamlConfig(command.CmdConfFile)
	//初始化BDB微型数据库
	if !utils.PathExists(path.Dir(common.Conf.BDB.Path)) {
		_ = os.MkdirAll(path.Dir(common.Conf.BDB.Path), 0775)
	}
	common.BDB, err = storm.Open(common.Conf.BDB.Path)
	if err != nil {
		fmt.Println("open storm database error:", err)
	}
	//初始化mongo db 连接池
	//err = ckdb.InitMongo(common.Conf.MDB)
	//if err != nil {
	//	fmt.Println("init mongo error:", err)
	//	os.Exit(1)
	//	return
	//}
	//初始化 redis 连接池
	//components.InitRedisPool(common.Conf.RDB)
	//写入PID文件
	if common.Conf.System.Pid != "" {
		command.CmdPidName = common.Conf.System.Pid
	}
	utils.WritePid(command.CmdPidName)
	//初始化关闭信号
	sigs = make(chan os.Signal, 1)
	done = make(chan bool, 1)
	//初始化全局内存缓存
	common.MemCache = components.NewMemCache()
	if command.CmdServer {
		//初始化HTTP WEB服务
		httpServer = router.NewHttpServer(common.Conf.System.Ip+":"+common.Conf.System.Port, command.CmdDebug, command.CmdCross, command.CmdPProf)
		httpServer.StaticEmbedFS(htmlFiles)
		//初始化TCP 主服务
		service.MainServer = service.NewTcpServer(fmt.Sprintf("%s:%s", common.Conf.Server.Ip, common.Conf.Server.Port), command.CmdDebug)
		common.SocketIO = websocket.NewEngine()
	} else if command.CmdNode {
		//初始化TCP 节点服务
		service.NodeServer = service.NewNodeService(common.Conf.Node.Server, common.Conf.Node.Name, common.Conf.Node.AuthPass)
	}
}

func Version() {
	fmt.Println("App Name:", AppName)
	fmt.Println("App Version:", AppVersion)
	fmt.Println("Build Version:", BuildVersion)
	fmt.Println("Build Time:", BuildTime)
	fmt.Println("Git Revision:", GitRevision)
	fmt.Println("Git Branch:", GitBranch)
	fmt.Println("Golang Version:", GoVersion)
}
