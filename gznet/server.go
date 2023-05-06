package gznet

import (
	"fmt"
	"net"

	"github.com/marsxingzhi/gozinx/gzinterface"
	"github.com/marsxingzhi/gozinx/handler"
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
	MsgHandler handler.IMsghandler
}

func New(name, ip string, port int) gzinterface.IServer {
	fmt.Println("server.New...")
	return &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         ip,
		Port:       port,
		MsgHandler: handler.New(),
	}
}

// 由于使用了Router，处理逻辑交给Router
// var handle gzinterface.HandleFunc = func(tcpConn *net.TCPConn, buf []byte, cnt int) error {
// 	newBuf := buf[:cnt]
// 	newBuf = append(newBuf, ", and back from server"...)
// 	_, err := tcpConn.Write(newBuf)
// 	if err != nil {
// 		fmt.Println("failed to write back to connection: ", err)
// 	}
// 	return err
// }

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

	fmt.Printf("start gozinx server %s success, and now listenning...\n", s.Name)

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

		// 4. 从connection中读取客户端传来的数据
		// @xingzhi 思考为什么处理业务不应该放在这个for循环中，而是另开一个goroutine
		// go handleConnection(conn)

		c := NewConnection(conn, connID, s.MsgHandler)
		connID++
		go c.Start()
	}

}

// 迁移到Connection中
// func handleConnection(conn *net.TCPConn) {
// 	defer conn.Close()
// 	buf := make([]byte, 1024)
// 	cnt, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("[server] failed to read from connection: ", err)
// 		return
// 	}

// 	fmt.Printf("[server] read data from connection: %s\n", string(buf))

// 	// test 回写到端上
// 	_, err = conn.Write(buf[:cnt])
// 	if err != nil {
// 		fmt.Println("[server] failed to write back to connection: ", err)
// 		return
// 	}
// }

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	go s.Start()

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, r gzinterface.IRouter) {
	s.MsgHandler.AddRouter(msgID, r)
}
