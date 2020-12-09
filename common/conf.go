package common

import (
	"github.com/clakeboy/golib/ckdb"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//总配置结构
type Config struct {
	System    *SystemConfig    `json:"system" yaml:"system"`
	DB        *ckdb.DBConfig   `json:"db" yaml:"db"`
	BDB       *BoltDBConfig    `json:"boltdb" yaml:"boltdb"`
	Cookie    *CookieConfig    `json:"cookie" yaml:"cookie"`
	HttpProxy *HttpProxyConfig `json:"http_proxy" yaml:"http_proxy"`
	Queue     *QueueConfig     `json:"queue" yaml:"queue"`
}

type SystemConfig struct {
	Port string `json:"port" yaml:"port"`
	Ip   string `json:"ip" yaml:"ip"`
	Pid  string `json:"pid" yaml:"pid"`
}

//Cookie 配置
type CookieConfig struct {
	Path     string `json:"path" yaml:"path"`
	Domain   string `json:"domain" yaml:"domain"`
	Source   bool   `json:"source" yaml:"source"`
	HttpOnly bool   `json:"http_only" yaml:"http_only"`
}

//boltdb 配置
type BoltDBConfig struct {
	Path string `json:"path" yaml:"path"`
}

//HTTP代理信息
type HttpProxyConfig struct {
	Use  bool   `json:"use" yaml:"use"`   //是否使用代理
	Addr string `json:"addr" yaml:"addr"` //HTTP代理地址
}

//队列配置
type QueueConfig struct {
	PushQueue  *QueueInfo `json:"push_queue" yaml:"push_queue"`
	UsersQueue *QueueInfo `json:"users_queue" yaml:"users_queue"`
}

//队列单个信息
type QueueInfo struct {
	Name   string `json:"name" yaml:"name"`
	Number int    `json:"number" yaml:"number"`
}

//节点配置
type NodeInfo struct {
}

//读取一个YAML配置文件
func NewYamlConfig(confFile string) *Config {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		panic(err)
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
