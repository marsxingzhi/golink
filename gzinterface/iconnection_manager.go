package gzinterface

// 链接管理模块的接口
type IConnectionManager interface {
	// 添加链接
	Add(conn IConnection)
	// 删除链接
	Remove(conn IConnection)
	// 获取链接
	Get(connID uint32) (IConnection, bool)
	// 链接个数
	Len() int
	// 清理所有链接
	ClearAllConns()
}
