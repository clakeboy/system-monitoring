package common

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/components"
	"system-monitoring/command"
)

var Conf *Config
var BDB *storm.DB
var MDB *ckdb.DBMongo
var MemCache *components.MemCache

func DebugF(str string, args ...interface{}) {
	if command.CmdDebug {
		fmt.Printf(str+"\n", args...)
	}
}
