package models

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/utils"
	"os"
	"path"
	"strings"
	"system-monitoring/common"
	"testing"
	"time"
)

func init() {
	var err error
	//获取YAML
	common.Conf = common.NewYamlConfig("../dev.conf")
	//初始化BDB微型数据库
	if !utils.PathExists(path.Dir(common.Conf.BDB.Path)) {
		_ = os.MkdirAll(path.Dir(common.Conf.BDB.Path), 0775)
	}
	common.BDB, err = storm.Open("../db/sys.db")
	if err != nil {
		fmt.Println("open storm database error:", err)
	}
}

func TestNewManagerModel(t *testing.T) {
	model := NewManagerModel(nil)
	data := new(ManagerData)
	data.CreatedDate = time.Now().Unix()
	data.Name = "admin"
	data.Password = utils.EncodeMD5("123123")
	data.Account = "admin"
	_ = model.Save(data)
}

func TestPath(t *testing.T) {
	name := "text.txt"
	name2 := "clake"
	ext := path.Ext(name)
	strings.ReplaceAll(name, ext, "")
	fmt.Println(path.Ext(name), path.Base(name))
	fmt.Println(path.Ext(name2), path.Base(name2))
}
