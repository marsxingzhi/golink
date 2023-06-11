package server

import (
	conn "github.com/marsxingzhi/xzlink/pkg/connection"
	"github.com/marsxingzhi/xzlink/pkg/router"
)

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 添加路由
	AddRouter(msgID uint32, r router.IRouter)
	// 获取链接管理模块
	GetConnectionManager() conn.IConnectionManager

	// hook函数
	SetOnConnStart(hookFunc func(conn conn.IConnection))
	SetOnConnStop(hookFunc func(conn conn.IConnection))
	CallOnConnStart(conn conn.IConnection)
	CallOnConnStop(conn conn.IConnection)
}
