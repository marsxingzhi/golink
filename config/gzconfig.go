package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/marsxingzhi/gozinx/gzinterface"
)

// 全局配置
type Gzconfig struct {
	// 当前服务器名称
	Name string `json:"name"`
	// 当前全局的server对象
	TcpServer gzinterface.IServer
	// 主机ip
	Host string `json:"host"`
	// 端口号
	Port int `json:"port"`

	// gozinx框架的版本号
	Version string `json:"version"`
	// 最大连接数
	MaxConn int `json:"max_connection"`
	// 最大数据包
	MaxPackageSize int32 `json:"max_package_size"`
}

var Config *Gzconfig

// 初始化配置
func Init() {
	// 默认配置
	defaultConfig()
	// 加载配置
	loadConfig()
}

func loadConfig() {
	// TODO 先写死
	bytes, err := ioutil.ReadFile("conf/gzinx.json")
	fmt.Printf("loadConfig data: %s\n", string(bytes))

	if err != nil {
		panic("[gozinx] failed to load config")
	}
	if err = json.Unmarshal(bytes, &Config); err != nil {
		fmt.Printf("[gozinx] failed to unmarshal gzinx.json: %+v\n", err)
		panic("[gozinx] failed to unmarshal gzinx.json")
	}
	fmt.Printf("[gozinx] load config success, and config: %+v\n", Config)
}

func defaultConfig() {
	Config = &Gzconfig{
		Name:           "gozinx server app",
		Host:           "0.0.0.0",
		Port:           8081,
		Version:        "0.1",
		MaxConn:        1000,
		MaxPackageSize: 1024,
	}
}
