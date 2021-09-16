package common

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/components"
	"system-monitoring/command"
)

var Conf *Config
var BDB *storm.DB
var MemCache *components.MemCache
var dbs map[string]*storm.DB

var debugLog = components.NewSysLog("debug_")

func DebugF(str string, args ...interface{}) {
	if command.CmdDebug {
		fmt.Printf(str+"\n", args...)
	} else {
		debugLog.Info("[DEBUG] " + fmt.Sprintf(str+"\n", args...))
	}
}

func GetNodeDB(name string) {

}
