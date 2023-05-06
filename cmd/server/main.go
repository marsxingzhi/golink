package main

import (
	"fmt"

	"github.com/marsxingzhi/golink/config"
	"github.com/marsxingzhi/golink/gzinterface"
	"github.com/marsxingzhi/golink/gznet"
)

type PingRouter struct {
	gznet.BaseRouter
}

func (pr *PingRouter) PreHandle(req gzinterface.IRequest) {
	// fmt.Println("call PreHandle")
}

func (pr *PingRouter) Handle(req gzinterface.IRequest) {
	fmt.Println("call PingRouter...")

	// req.GetConnection().GetConn().Write([]byte("ping...\n"))

	// 先读取客户端的数据，再回写ping数据
	fmt.Printf("[Server] receive msg | msgID: %v, dataLen: %v, data: %v\n", req.GetMsgID(), len(req.GetData()), string(req.GetData()))

	if err := req.GetConnection().SendMessage(1000, []byte("ping...\n")); err != nil {
		fmt.Printf("[Server] | failed to sendMessage: %v\n", err)
		return
	}
}

func (pr *PingRouter) PostHandle(req gzinterface.IRequest) {
	// fmt.Println("call PostHandle")
}

type HelloRouter struct {
	gznet.BaseRouter
}

func (pr *HelloRouter) Handle(req gzinterface.IRequest) {
	fmt.Println("call HelloRouter...")

	// req.GetConnection().GetConn().Write([]byte("ping...\n"))

	// 先读取客户端的数据，再回写ping数据
	fmt.Printf("[Server] receive msg | msgID: %v, dataLen: %v, data: %v\n", req.GetMsgID(), len(req.GetData()), string(req.GetData()))

	if err := req.GetConnection().SendMessage(1001, []byte("hello...\n")); err != nil {
		fmt.Printf("[Server] | failed to sendMessage: %v\n", err)
		return
	}
}

func onConnStart(conn gzinterface.IConnection) {
	fmt.Println("main onConnStart...")
	// 这里也可以提前发一下消息
}

func onConnStop(conn gzinterface.IConnection) {
	fmt.Println("main onConnStop...")
}

func main() {

	// 初始化配置
	config.Init()

	// 1. 创建server
	server := gznet.New("t1", "127.0.0.1", 8081)

	// 2. 添加router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})

	// 注册
	server.SetOnConnStart(onConnStart)
	server.SetOnConnStop(onConnStop)

	// 3. 启动server
	server.Serve()
}
