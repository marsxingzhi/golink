package datapack

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/marsxingzhi/golink/config"
	"github.com/marsxingzhi/golink/model"
)

func TestDataPack(t *testing.T) {
	config.Init()

	listener, err := net.Listen("tcp", "127.0.0.1:8100")
	if err != nil {
		fmt.Println("failed to server listen: ", err)
		return
	}

	// 服务端读取数据
	go func() {

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("failed to accept: ", err)
				continue
			}

			// 读取数据
			go func(conn net.Conn) {
				dp := New()
				for {
					headData := make([]byte, dp.GetHeadLen())

					// 读满这个buf
					if _, err = io.ReadFull(conn, headData); err != nil {
						fmt.Println("failed to read headData: ", err)
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("failed to UnPack: ", err)
						return
					}
					if msgHead.DataLen > 0 {
						// 开始读取消息内容
						msg := msgHead
						msg.Data = make([]byte, msg.DataLen)

						if _, err = io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("failed to readfull msg.Data: ", err)
							return
						}

						fmt.Printf("receive message ID: %v, DataLen: %v, Data: %v\n", msg.MsgID, msg.DataLen, string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	go func() {
		// 客户端
		conn, err := net.Dial("tcp", "127.0.0.1:8100")
		if err != nil {
			fmt.Println("failed to dial: ", err)
			return
		}

		dp := New()

		// 模拟粘包过程，封装两个message，一起发送
		msg1 := &model.Message{
			MsgID:   1,
			DataLen: 5,
			Data:    []byte{'1', '2', '3', '4', '5'},
		}
		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("failed to pack msg1: ", err)
			return
		}

		msg2 := &model.Message{
			MsgID:   2,
			DataLen: 6,
			Data:    []byte{'a', 'b', 'c', 'd', 'e', 'f'},
		}

		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("failed to pack msg2: ", err)
			return
		}

		res := make([]byte, 0)
		res = append(res, sendData1...)
		res = append(res, sendData2...)

		fmt.Printf("res: %+v\n", res)
		fmt.Println("sendData1 len: ", len(sendData1))
		fmt.Println("sendData2 len: ", len(sendData2))

		if _, err = conn.Write(res); err != nil {
			fmt.Println("failed to write sendData: ", err)
			return
		}
	}()

	// 因为需要看到服务端的处理结果，这里阻塞一下
	select {}

}
