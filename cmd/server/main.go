package main

import (
	"fmt"

	"github.com/marsxingzhi/gozinx/config"
	"github.com/marsxingzhi/gozinx/gzinterface"
	"github.com/marsxingzhi/gozinx/gznet"
)

type PingRouter struct {
	gznet.BaseRouter
}

func (pr *PingRouter) PreHandle(req gzinterface.IRequest) {
	// fmt.Println("call PreHandle")
}

func (pr *PingRouter) Handle(req gzinterface.IRequest) {
	fmt.Println("call Handle...")

	// req.GetConnection().GetConn().Write([]byte("ping...\n"))

	// 先读取客户端的数据，再回写ping数据
	fmt.Printf("[Server] receive msg | msgID: %v, dataLen: %v, data: %v\n", req.GetMsgID(), len(req.GetData()), string(req.GetData()))

	if err := req.GetConnection().SendMessage(1, []byte("ping...\n")); err != nil {
		fmt.Printf("[Server] | failed to sendMessage: %v\n", err)
		return
	}
}

func (pr *PingRouter) PostHandle(req gzinterface.IRequest) {
	// fmt.Println("call PostHandle")
}

func main() {

	// 初始化配置
	config.Init()

	// 1. 创建server
	server := gznet.New("t1", "127.0.0.1", 8081)

	// 2. 添加router
	server.AddRouter(&PingRouter{})

	// 3. 启动server
	server.Serve()
}
