package common

import (
	"github.com/clakeboy/golib/ckdb"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//总配置结构
type Config struct {
	path      string           `json:"-" yaml:"-"`
	System    *SystemConfig    `json:"system" yaml:"system"`
	DB        *ckdb.DBConfig   `json:"db" yaml:"db"`
	BDB       *BoltDBConfig    `json:"boltdb" yaml:"boltdb"`
	Cookie    *CookieConfig    `json:"cookie" yaml:"cookie"`
	HttpProxy *HttpProxyConfig `json:"http_proxy" yaml:"http_proxy"`
	Server    *Server          `json:"server" yaml:"server"` //主服务模式配置
	Node      *NodeInfo        `json:"node" yaml:"node"`     //节点服务模式配置
}

//HTTP及系统配置
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

//节点服务配置
type NodeInfo struct {
	Server   string `json:"server" yaml:"server"`       //服务器地址
	Name     string `json:"name" yaml:"name"`           //节点名称
	AuthPass string `json:"auth_pass" yaml:"auth_pass"` //认证密钥
}

//主服务模式配置
type Server struct {
	Ip       string `json:"ip" yaml:"ip"`               //指定服务IP,为空就是接受所有IP
	Port     string `json:"port" yaml:"port"`           //指定服务端口,默认为 17711
	AuthPass string `json:"auth_pass" yaml:"auth_pass"` //认证密钥
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

//写入YAML配置
func (c *Config) Save() {
	out, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	ioutil.WriteFile("./run.conf", out, 0755)
}
