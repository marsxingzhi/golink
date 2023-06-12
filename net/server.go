package net

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/pkg/config"
	conn "github.com/marsxingzhi/xzlink/pkg/connection"
	"github.com/marsxingzhi/xzlink/pkg/msg_handler"
	"github.com/marsxingzhi/xzlink/pkg/router"
	"github.com/marsxingzhi/xzlink/pkg/server"
	"net"
)

// IServer接口的实现
type Server struct {
	// 服务器名称
	Name string
	// IP版本
	IPVersion string
	// 服务器绑定的ip地址
	IP string
	// 服务器绑定的端口
	Port int

	// Router gzinterface.IRouter
	MsgHandler msg_handler.IMsgHandler

	// 链接管理
	ConnMgr conn.IConnectionManager

	// hook函数
	// 链接创建之后
	OnConnStart func(conn conn.IConnection)
	// 链接关闭之前
	OnConnStop func(conn conn.IConnection)
}

func NewServer(name, ip string, port int) server.IServer {
	fmt.Println("server.New...")
	return &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         ip,
		Port:       port,
		MsgHandler: msg_handler.New(),
		ConnMgr:    NewConnectionManager(),
	}
}

func (s *Server) Start() {
	fmt.Println("server.Start...")

	// 0. 开启工作池
	s.MsgHandler.StartWorkerPool()

	// 1. 获取TCPAddr
	tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("failed to ResolveTCPAddr: ", err)
		return
	}
	// 2. 根据addr，拿到listener
	listener, err := net.ListenTCP(s.IPVersion, tcpAddr)
	if err != nil {
		fmt.Println("failed to ListeneTCP: ", err)
		return
	}

	fmt.Printf("start xzlink server %s success, and now listenning...\n", s.Name)

	var connID uint32 = 0

	// @marsxingzhi  这个循环不能随便退出。下面的c.Start方法可以结束
	for {
		// 3. accept，拿到connection
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept connection: ", err)
			continue
		}
		fmt.Printf("new connection from %s\n", conn.RemoteAddr())

		fmt.Println("[Server] Len: ", s.ConnMgr.Len(), ", MaxConn: ", config.Config.GetMaxConn())
		// 判断链接是否超过最大值
		if s.ConnMgr.Len() >= config.Config.GetMaxConn() {
			conn.Close()
			fmt.Println("[Server] too mutch connections...")
			continue
		}

		// 4. 从connection中读取客户端传来的数据
		// @xingzhi 思考为什么处理业务不应该放在这个for循环中，而是另开一个goroutine
		// go handleConnection(conn)

		c := NewConnection(s, conn, connID, s.MsgHandler)
		connID++
		go c.Start()
	}

}

func (s *Server) Stop() {
	// 清理资源
	s.ConnMgr.ClearAllConns()

}

func (s *Server) Serve() {
	go s.Start()

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, r router.IRouter) {
	s.MsgHandler.AddRouter(msgID, r)
}

func (s *Server) GetConnectionManager() conn.IConnectionManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(conn conn.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(conn conn.IConnection)) {
	s.OnConnStop = hookFunc
}
func (s *Server) CallOnConnStart(conn conn.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn conn.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
