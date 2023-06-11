package router

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/gznet"
	"github.com/marsxingzhi/xzlink/pkg/request"
)

type HelloRouter struct {
	gznet.BaseRouter
}

func (pr *HelloRouter) Handle(req request.IRequest) {
	fmt.Println("call HelloRouter...")

	// req.GetConnection().GetConn().Write([]byte("ping...\n"))

	// 先读取客户端的数据，再回写ping数据
	fmt.Printf("[Server] receive msg | msgID: %v, dataLen: %v, data: %v\n", req.GetMsgID(), len(req.GetData()), string(req.GetData()))

	if err := req.GetConnection().SendMessage(1001, []byte("hello...\n")); err != nil {
		fmt.Printf("[Server] | failed to sendMessage: %v\n", err)
		return
	}
}
