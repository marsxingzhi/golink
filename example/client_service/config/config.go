package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type XzConfig struct {
	Client ClientConfig `yaml:"client"`
}

// ClientConfig 全局配置
type ClientConfig struct {
	// 最大数据包
	MaxPackageSize int32 `yaml:"max_package_size"`
}

var Config *XzConfig

func (gc *XzConfig) GetMaxPackageSize() int32 {
	return gc.Client.MaxPackageSize
}

// Init 初始化配置
func Init(path string) {
	// 加载配置
	Config = new(XzConfig)
	loadConfig(path)
}

func loadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	fmt.Printf("load client Config data: %s\n", string(bytes))

	if err != nil {
		fmt.Printf("[xzlink] failed to load client config: %+v\n", err)
		panic("[xzlink] failed to load config")
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		fmt.Printf("[xzlink] failed to unmarshal xzlink.yaml: %+v\n", err)
		panic("[xzlink] failed to unmarshal xzlink.yaml")
	}

	fmt.Printf("[xzlink] load client config success, and config: %+v\n", Config)
}
