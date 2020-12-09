package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"os"
	"path"
	"system-monitoring/command"
	"system-monitoring/common"
	"system-monitoring/service"
)

var out chan os.Signal
var server *service.HttpServer

var (
	AppName      string //应用名称
	AppVersion   string //应用名称
	BuildVersion string //编译版本
	BuildTime    string //编译时间
	GitRevision  string //Git 版本
	GitBranch    string //Git 分支
	GoVersion    string //Golang 信息
)

func main() {
	go utils.ExitApp(out, func(s os.Signal) {
		_ = os.Remove(command.CmdPidName)
	})
	server.Start()
}

func init() {
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
	out = make(chan os.Signal, 1)
	//初始化全局内存缓存
	common.MemCache = components.NewMemCache()
	//初始化HTTP WEB服务
	server = service.NewHttpServer(common.Conf.System.Ip+":"+common.Conf.System.Port, command.CmdDebug, command.CmdCross, command.CmdPProf)
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
