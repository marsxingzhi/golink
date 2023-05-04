package gznet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/marsxingzhi/gozinx/datapack"
	"github.com/marsxingzhi/gozinx/gzinterface"
	"github.com/marsxingzhi/gozinx/model"
)

// conn对象
type Connection struct {
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态，是否已经关闭
	IsClose bool
	// 与当前链接所绑定的处理业务的函数
	// Handle gzinterface.HandleFunc
	// 等待链接退出的channel
	ExitChan chan []byte
	Router   gzinterface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, r gzinterface.IRouter) *Connection {
	return &Connection{
		Conn:    conn,
		ConnID:  connID,
		IsClose: false,
		// Handle:   handle,
		ExitChan: make(chan []byte, 1),
		Router:   r,
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
		// buf := make([]byte, 1024)
		// cnt, err := c.Conn.Read(buf)
		// if err != nil {
		// 	if errors.Is(err, io.EOF) {
		// 		fmt.Println("end of data")
		// 		break
		// 	}
		// 	fmt.Println("failed to read from connection: ", err)
		// 	continue
		// }
		// fmt.Printf("read from connection success, and msg: %s\n", string(buf[:cnt]))

		// 上述代码注释掉，这里使用拆包的方式
		dp := datapack.New()
		// 1. 先读取head，获取到消息长度
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Printf("[Connection] Start | failed to read msg head: %v\n", err)
			break
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Printf("[Connection] Satrt | failed to UnPack: %v\n", err)
			break
		}

		// 2. 根据消息长度，读取消息内容
		if msg.DataLen > 0 {
			msg.Data = make([]byte, msg.DataLen)

			if _, err = io.ReadFull(c.Conn, msg.Data); err != nil {
				fmt.Printf("[Connection] Start | failed to readfull msg data: %v\n", err)
				break
			}
		}

		// 交给Router处理
		// if err = c.Handle(c.Conn, buf, cnt); err != nil {
		// 	fmt.Printf("ConnID %v handle is error\n", c.ConnID)
		// 	break
		// }

		// 这里有必要开goroutine？
		req := Request{
			Conn: c,
			Msg:  msg,
		}
		c.Router.PreHandle(&req)
		c.Router.Handle(&req)
		c.Router.PostHandle(&req)

	}

}

func (c *Connection) Stop() {

}

func (c *Connection) GetConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMessage 先封包，再发送
func (c *Connection) SendMessage(msgID uint32, data []byte) error {
	if c.IsClose {
		fmt.Println("[Connection] SendMessage | connection has be closed")
		return errors.New("[Connection] SendMessage | connection has be closed")
	}

	dp := datapack.New()

	msg := &model.Message{
		MsgID:   msgID,
		DataLen: uint32(len(data)),
		Data:    data,
	}

	sendData, err := dp.Pack(msg)
	if err != nil {
		fmt.Printf("[Connection] SendMessage | failed to Pack msg: %v\n", err)
		return err
	}
	// 将数据写回客户端
	if _, err = c.Conn.Write(sendData); err != nil {
		fmt.Printf("[Connection] SendMessage | failed to Write msg: %v\n", err)
		return err
	}
	return nil
}
