package gznet

import (
	"fmt"
	"net"

	"github.com/marsxingzhi/gozinx/gzinterface"
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
}

func New(name, ip string, port int) gzinterface.IServer {
	fmt.Println("server.New...")
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        ip,
		Port:      port,
	}
}

func (s *Server) Start() {
	fmt.Println("[server] server.Start...")
	// 1. 获取TCPAddr
	tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("[server] failed to ResolveTCPAddr: ", err)
		return
	}
	// 2. 根据addr，拿到listener
	listener, err := net.ListenTCP(s.IPVersion, tcpAddr)
	if err != nil {
		fmt.Println("[server] failed to ListeneTCP: ", err)
		return
	}

	fmt.Printf("[server] start gozinx server %s success, and now listenning...\n", s.Name)

	for {
		// 3. accept，拿到connection
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("[server] failed to accept connection: ", err)
			continue
		}
		fmt.Printf("[server] new connection from %s\n", conn.RemoteAddr())

		// 4. 从connection中读取客户端传来的数据
		// @xingzhi 思考为什么处理业务不应该放在这个for循环中，而是另开一个goroutine
		go handleConnection(conn)
	}

}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println("[server] failed to read from connection: ", err)
		return
	}

	fmt.Printf("[server] read data from connection: %s\n", string(buf))

	// test 回写到端上
	_, err = conn.Write(buf[:cnt])
	if err != nil {
		fmt.Println("[server] failed to write back to connection: ", err)
		return
	}
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	go s.Start()

	// 阻塞
	select {}
}
