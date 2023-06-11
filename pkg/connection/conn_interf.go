package connection

import "net"

// connection的接口
type IConnection interface {
	// 启动链接
	Start()
	// 停止链接
	Stop()
	// 获取当前链接的conn对象
	GetConn() *net.TCPConn
	// 获取链接ID
	GetConnID() uint32
	// 获取客户端链接的地址和端口
	RemoteAddr() net.Addr
	// 发送数据
	SendMessage(msgID uint32, data []byte) error

	// 链接属性
	SetProperty(key string, value interface{})
	GetProperity(key string) (interface{}, bool)
	RemoveProperity(key string)
}

// 函数类型
// 参数一：conn对象；参数二：处理的数据内容；参数三：处理的数据量
type HandleFunc func(*net.TCPConn, []byte, int) error
