package main

import (
	"fmt"

	"github.com/marsxingzhi/gozinx/gznet"
)

func main() {
	fmt.Println("start server")

	// 1. 创建server
	server := gznet.New("t1", "127.0.0.1", 8081)
	// 2. 启动server
	server.Serve()
}
