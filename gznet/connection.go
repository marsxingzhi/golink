package gznet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/marsxingzhi/golink/datapack"
	"github.com/marsxingzhi/golink/gzinterface"
	"github.com/marsxingzhi/golink/handler"
	"github.com/marsxingzhi/golink/model"
)

// conn对象
type Connection struct {
	// 当前Connection属于哪个Server
	TcpServer gzinterface.IServer

	// TCP的链接
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态，是否已经关闭
	IsClose bool
	// 与当前链接所绑定的处理业务的函数
	// Handle gzinterface.HandleFunc
	// 等待链接退出的channel
	ExitChan chan bool

	MsgHandler handler.IMsghandler

	// 无缓冲，读写分离，Reader与Writer通信
	MsgChan chan []byte

	// 链接属性
	Properity      map[string]interface{}
	ProperityMutex sync.RWMutex
}

func NewConnection(server gzinterface.IServer, conn *net.TCPConn, connID uint32, msgHandler handler.IMsghandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClose:    false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		MsgChan:    make(chan []byte),
		Properity:  make(map[string]interface{}),
	}

	// 将Connection添加到ConnMgr中
	c.TcpServer.GetConnectionManager().Add(c)

	return c
}

// 需要开启goroutine
func (c *Connection) Start() {
	go c.StartRead()
	go c.StartWrite()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.IsClose {
		return
	}

	c.IsClose = true
	// 在连接关闭之前调用
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	// 回收资源
	close(c.MsgChan)
	close(c.ExitChan)

	// 删除链接
	c.TcpServer.GetConnectionManager().Remove(c)
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
	// if _, err = c.Conn.Write(sendData); err != nil {
	// 	fmt.Printf("[Connection] SendMessage | failed to Write msg: %v\n", err)
	// 	return err
	// }

	// 这里就不是直接将数据写回客户端，而是写到channel中，由专门的Writer负责回写
	c.MsgChan <- sendData

	return nil
}

func (c *Connection) StartRead() {
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

	fmt.Println("[Connection] StartRead...")

	// defer fmt.Printf("connection closed, and ConnID: %v, remote addr is %s\n", c.ConnID, c.RemoteAddr())
	// defer c.Stop()

	defer func() {
		c.Stop()
		fmt.Printf("[Connection] closed, and ConnID: %v, remote addr is %s\n", c.ConnID, c.RemoteAddr())
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
			fmt.Printf("[Connection] StartRead | failed to read msg head: %v\n", err)
			break
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Printf("[Connection] StartRead | failed to UnPack: %v\n", err)
			break
		}

		// 2. 根据消息长度，读取消息内容
		if msg.DataLen > 0 {
			msg.Data = make([]byte, msg.DataLen)

			if _, err = io.ReadFull(c.Conn, msg.Data); err != nil {
				fmt.Printf("[Connection] StartRead | failed to readfull msg data: %v\n", err)
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
		// c.Router.PreHandle(&req)
		// c.Router.Handle(&req)
		// c.Router.PostHandle(&req)

		// go c.MsgHandler.DoHandle(&req) // 可以开个goroutine，将处理request的逻辑go出去，然后继续读新的数据

		// 这里必须要用go，方法体内只是往消息队列发个消息
		c.MsgHandler.SendMessageToTaskQueue(&req)

	}
}

// StartWrite 是Writer往客户端写数据的逻辑。需要go
// 不要传参数，用channel可以解耦
func (c *Connection) StartWrite() {
	fmt.Println("[Connection] StartWrite...")

	for {
		select {
		case data := <-c.MsgChan: // 收到Reader给的数据
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("[Connection] StartWrite | failed to Write msg back to client: %v\n", err)
			}

		case <-c.ExitChan: // 收到Reader给的退出信号
			return
		}
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.ProperityMutex.Lock()
	defer c.ProperityMutex.Unlock()
	c.Properity[key] = value
}
func (c *Connection) GetProperity(key string) (interface{}, bool) {
	c.ProperityMutex.RLock()
	defer c.ProperityMutex.RUnlock()

	if val, ok := c.Properity[key]; ok {
		return val, true
	} else {
		return nil, false
	}
}
func (c *Connection) RemoveProperity(key string) {
	c.ProperityMutex.Lock()
	defer c.ProperityMutex.Unlock()
	delete(c.Properity, key)
}
