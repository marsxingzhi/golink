package router

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/net"
	"github.com/marsxingzhi/xzlink/pkg/request"
)

type PingRouter struct {
	net.BaseRouter
}

func (pr *PingRouter) PreHandle(req request.IRequest) {
	// fmt.Println("call PreHandle")
}

func (pr *PingRouter) Handle(req request.IRequest) {
	fmt.Println("call PingRouter...")

	// req.GetConnection().GetConn().Write([]byte("ping...\n"))

	// 先读取客户端的数据，再回写ping数据
	fmt.Printf("[Server] receive msg | msgID: %v, dataLen: %v, data: %v\n", req.GetMsgID(), len(req.GetData()), string(req.GetData()))

	if err := req.GetConnection().SendMessage(1000, []byte("ping...\n")); err != nil {
		fmt.Printf("[Server] | failed to sendMessage: %v\n", err)
		return
	}
}

func (pr *PingRouter) PostHandle(req request.IRequest) {
	// fmt.Println("call PostHandle")
}
