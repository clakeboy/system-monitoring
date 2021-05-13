package command

import (
	"flag"
	"fmt"
	"os"
)

var (
	CmdDebug       bool
	CmdCross       bool
	CmdPProf       bool
	CmdConfFile    string
	CmdPidName     string
	CmdShowVersion bool
	CmdDaemon      bool //是否启动守护进程
	CmdServer      bool //单主服务模式运行
	CmdNode        bool //单节点模式运行
	CmdPassive     bool //节点被动模式
)

func InitCommand() {
	flag.BoolVar(&CmdDebug, "debug", false, "is runtime debug mode")
	flag.BoolVar(&CmdCross, "cross", false, "use cross request")
	flag.BoolVar(&CmdPProf, "pprof", false, "open go pprof debug")
	flag.StringVar(&CmdConfFile, "config", "./main.conf", "app config file")
	flag.StringVar(&CmdPidName, "pid", "./monitoring.pid", "app config file")
	flag.BoolVar(&CmdShowVersion, "version", false, "show this version information")
	flag.BoolVar(&CmdDaemon, "daemon", false, "start daemon")
	flag.BoolVar(&CmdServer, "server", false, "start only server mode")
	flag.BoolVar(&CmdNode, "node", false, "start only node mode")
	flag.BoolVar(&CmdPassive, "passive", false, "node passive connect for server")
	flag.Parse()
	ExecCommand()
}

func ExecCommand() {
	if CmdDaemon {
		StartDaemon()
	}
}

//结束程序
func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
