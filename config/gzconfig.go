package config

import (
	"fmt"
	"io/ioutil"

	"github.com/marsxingzhi/gozinx/gzinterface"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml: server`
}

// 全局配置
type ServerConfig struct {
	// 当前服务器名称
	Name string `yaml:"name"`
	// 当前全局的server对象
	TcpServer gzinterface.IServer
	// 主机ip
	Host string `yaml:"host"`
	// 端口号
	Port int `yaml:"port"`

	// gozinx框架的版本号
	Version string `yaml:"version"`
	// 最大连接数
	MaxConn int `yaml:"max_connection"`
	// 最大数据包
	MaxPackageSize int32 `yaml:"max_package_size"`
}

var GzConfig *Config

func (gc *Config) GetMaxPackageSize() int32 {
	return gc.Server.MaxPackageSize
}

// 初始化配置
func Init() {
	// 默认配置
	defaultConfig()
	// 加载配置
	// loadConfig()
}

func loadConfig() {
	// TODO 先写死
	bytes, err := ioutil.ReadFile("conf/config.yaml")
	fmt.Printf("loadConfig data: %s\n", string(bytes))

	if err != nil {
		panic("[gozinx] failed to load config")
	}
	if err = yaml.Unmarshal(bytes, &GzConfig); err != nil {
		fmt.Printf("[gozinx] failed to unmarshal gzinx.json: %+v\n", err)
		panic("[gozinx] failed to unmarshal gzinx.json")
	}

	fmt.Printf("[gozinx] load config success, and config: %+v\n", GzConfig)
}

func defaultConfig() {

	sc := ServerConfig{
		Name:           "gozinx server app",
		Host:           "0.0.0.0",
		Port:           8081,
		Version:        "0.1",
		MaxConn:        1000,
		MaxPackageSize: 1024,
	}

	GzConfig = &Config{
		Server: sc,
	}
}
