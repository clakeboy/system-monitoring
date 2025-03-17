package common

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/components/snowflake"
	"system-monitoring/command"
	"system-monitoring/websocket"
)

var Conf *Config
var BDB *storm.DB
var SocketIO *websocket.Engine
var MemCache *components.MemCache
var dbs map[string]*storm.DB
var SnowFlake *snowflake.SnowFlake

//var debugLog = components.NewSysLog("debug_")

func DebugF(str string, args ...interface{}) {
	if command.CmdDebug {
		fmt.Printf("[DEBUG] "+str+"\n", args...)
	}
}

func GetNodeDB(name string) {

}
