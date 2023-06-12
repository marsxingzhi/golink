package main

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/cmd/server/router"
	"github.com/marsxingzhi/xzlink/net"
	"github.com/marsxingzhi/xzlink/pkg/config"
	conn "github.com/marsxingzhi/xzlink/pkg/connection"
)

func onConnStart(conn conn.IConnection) {
	fmt.Println("main onConnStart...")
	// 这里也可以提前发一下消息
}

func onConnStop(conn conn.IConnection) {
	fmt.Println("main onConnStop...")
}

func main() {

	// 初始化配置
	configPath := "/Users/geyan/go/src/xzlink/configs/config.yaml"
	config.Init(configPath)

	// 1. 创建server
	server := net.NewServer("t1", "127.0.0.1", 8081)

	// 2. 添加router
	server.AddRouter(0, &router.PingRouter{})
	server.AddRouter(1, &router.HelloRouter{})

	// 注册
	server.SetOnConnStart(onConnStart)
	server.SetOnConnStop(onConnStop)

	// 3. 启动server
	server.Serve()
}
