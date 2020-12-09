package common

import (
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/components"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//总配置结构
type Config struct {
	System    *SystemConfig           `json:"system" yaml:"system"`
	DB        *ckdb.DBConfig          `json:"db" yaml:"db"`
	BDB       *BoltDBConfig           `json:"boltdb" yaml:"boltdb"`
	Cookie    *CookieConfig           `json:"cookie" yaml:"cookie"`
	MDB       *ckdb.MongoDBConfig     `json:"mdb" yaml:"mdb"`
	RDB       *components.RedisConfig `json:"redis" yaml:"redis"`
	Platform  *PlatConfig             `json:"platform" yaml:"platform"`
	HttpProxy *HttpProxyConfig        `json:"http_proxy" yaml:"http_proxy"`
	Queue     *QueueConfig            `json:"queue" yaml:"queue"`
	Oss       *OSSConfig              `json:"oss" yaml:"oss"`
	Up        *UploadConfig           `json:"upload" yaml:"upload"`
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

//微信开放平台信息
type PlatConfig struct {
	Token     string `json:"token" yaml:"token"`
	AesKey    string `json:"aes_key" yaml:"aes_key"`
	AppId     string `json:"appid" yaml:"appid"`
	AppSecret string `json:"app_secret" yaml:"app_secret"`
}

//OSS配置
type OSSConfig struct {
	Endpoint         string `json:"endpoint" yaml:"endpoint"`
	InternalEndpoint string `json:"internal" yaml:"internal"`
	AccessID         string `json:"accessID" yaml:"accessID"`
	AccessKey        string `json:"accessKey" yaml:"accessKey"`
	BucketName       string `json:"bucketName" yaml:"bucketName"`
	OssPath          string `json:"oss_path" yaml:"oss_path"`
	IsInternal       bool   `json:"is_internal" yaml:"is_internal"`
}

//文件上传配置
type UploadConfig struct {
	UploadMaxLength int    `json:"upload_max_length" yaml:"upload_max_length"`
	LocalSavePath   string `json:"local_save_path" yaml:"local_save_path"`
	DownloadDomain  string `json:"download_domain" yaml:"download_domain"`
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
