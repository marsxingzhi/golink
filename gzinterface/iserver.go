package gzinterface

// 服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 添加路由
	AddRouter(msgID uint32, r IRouter)
	// 获取链接管理模块
	GetConnectionManager() IConnectionManager

	// hook函数
	SetOnConnStart(hookFunc func(conn IConnection))
	SetOnConnStop(hookFunc func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}
