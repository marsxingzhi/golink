package config

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/pkg/server"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type XzConfig struct {
	Server ServerConfig `yaml:"XzConfig"`
}

// ServerConfig 全局配置
type ServerConfig struct {
	// 当前服务器名称
	Name string `yaml:"name"`
	// 当前全局的server对象
	TcpServer server.IServer
	// 主机ip
	Host string `yaml:"host"`
	// 端口号
	Port int `yaml:"port"`

	// xzlink框架的版本号
	Version string `yaml:"version"`
	// 最大连接数
	MaxConn int `yaml:"max_connection"`
	// 最大数据包
	MaxPackageSize int32 `yaml:"max_package_size"`
	// WorkerPool的大小
	WorkerPoolSize int `yaml:"worker_pool_size"`
	// 一个消息队列中的最大消息任务数
	MaxWorkerTaskCapacity int `yaml:"max_worker_task_capacity"`
}

var Config *XzConfig

func (gc *XzConfig) GetMaxPackageSize() int32 {
	return gc.Server.MaxPackageSize
}

func (gc *XzConfig) GetWorkerPoolSize() int {
	return gc.Server.WorkerPoolSize
}

func (gc *XzConfig) GetWorkerTaskCapacity() int {
	return gc.Server.MaxWorkerTaskCapacity
}

func (gc *XzConfig) GetMaxConn() int {
	return gc.Server.MaxConn
}

// Init 初始化配置
func Init(path string) {
	// 默认配置
	defaultConfig()
	// 加载配置
	loadConfig(path)
}

func loadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	fmt.Printf("loadConfig data: %s\n", string(bytes))

	if err != nil {
		fmt.Printf("[xzlink] failed to load config: %+v\n", err)
		panic("[xzlink] failed to load config")
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		fmt.Printf("[xzlink] failed to unmarshal gzinx.json: %+v\n", err)
		panic("[xzlink] failed to unmarshal xzlink.yaml")
	}

	fmt.Printf("[xzlink] load config success, and config: %+v\n", Config)
}

func defaultConfig() {

	sc := ServerConfig{
		Name:                  "xzlink server app",
		Host:                  "0.0.0.0",
		Port:                  80811,
		Version:               "0.1",
		MaxConn:               1000,
		MaxPackageSize:        1024,
		WorkerPoolSize:        10,
		MaxWorkerTaskCapacity: 1024,
	}

	Config = &XzConfig{
		Server: sc,
	}
}
