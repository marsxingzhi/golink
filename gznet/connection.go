package gznet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/marsxingzhi/gozinx/gzinterface"
)

// conn对象
type Connection struct {
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态，是否已经关闭
	IsClose bool
	// 与当前链接所绑定的处理业务的函数
	Handle gzinterface.HandleFunc
	// 等待链接退出的channel
	ExitChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, handle gzinterface.HandleFunc) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		IsClose:  false,
		Handle:   handle,
		ExitChan: make(chan []byte, 1),
	}
}

// 需要开启goroutine
func (c *Connection) Start() {
	// defer conn.Close()
	// buf := make([]byte, 1024)
	// cnt, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("[server] failed to read from connection: ", err)
	// 	return
	// }

	// fmt.Printf("[server] read data from connection: %s\n", string(buf))

	// // test 回写到端上
	// _, err = conn.Write(buf[:cnt])
	// if err != nil {
	// 	fmt.Println("[server] failed to write back to connection: ", err)
	// 	return
	// }

	fmt.Println("connection start...")

	// defer fmt.Printf("connection closed, and ConnID: %v, remote addr is %s\n", c.ConnID, c.RemoteAddr())
	// defer c.Stop()

	defer func() {
		c.Stop()
		fmt.Printf("connection closed, and ConnID: %v, remote addr is %s\n", c.ConnID, c.RemoteAddr())
	}()

	for {
		buf := make([]byte, 1024)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("end of data")
				break
			}
			fmt.Println("failed to read from connection: ", err)
			continue
		}
		fmt.Printf("read from connection success, and msg: %s\n", string(buf[:cnt]))

		if err = c.Handle(c.Conn, buf, cnt); err != nil {
			fmt.Printf("ConnID %v handle is error\n", c.ConnID)
			break
		}
	}

}

func (c *Connection) Stop() {

}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send() {

}
