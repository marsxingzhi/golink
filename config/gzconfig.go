package config

import (
	"fmt"
	"io/ioutil"

	"github.com/marsxingzhi/golink/gzinterface"
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

	// golink框架的版本号
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

var GzConfig *Config

func (gc *Config) GetMaxPackageSize() int32 {
	return gc.Server.MaxPackageSize
}

func (gc *Config) GetWorkerPoolSize() int {
	return gc.Server.WorkerPoolSize
}

func (gc *Config) GetWorkerTaskCapacity() int {
	return gc.Server.MaxWorkerTaskCapacity
}

func (gc *Config) GetMaxConn() int {
	return gc.Server.MaxConn
}

// 初始化配置
func Init() {
	// 默认配置
	defaultConfig()
	// 加载配置
	loadConfig()
}

func loadConfig() {
	// TODO 先写死
	bytes, err := ioutil.ReadFile("/Users/geyan/codes/golink/cmd/server/conf/config.yaml")
	fmt.Printf("loadConfig data: %s\n", string(bytes))

	if err != nil {
		fmt.Printf("[golink] failed to load config: %+v\n", err)
		panic("[golink] failed to load config")
	}
	if err = yaml.Unmarshal(bytes, &GzConfig); err != nil {
		fmt.Printf("[golink] failed to unmarshal gzinx.json: %+v\n", err)
		panic("[golink] failed to unmarshal gzinx.json")
	}

	fmt.Printf("[golink] load config success, and config: %+v\n", GzConfig)
}

func defaultConfig() {

	sc := ServerConfig{
		Name:                  "golink server app",
		Host:                  "0.0.0.0",
		Port:                  80811,
		Version:               "0.1",
		MaxConn:               1000,
		MaxPackageSize:        1024,
		WorkerPoolSize:        10,
		MaxWorkerTaskCapacity: 1024,
	}

	GzConfig = &Config{
		Server: sc,
	}
}
