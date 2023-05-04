package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/marsxingzhi/gozinx/config"
	"github.com/marsxingzhi/gozinx/datapack"
	"github.com/marsxingzhi/gozinx/model"
)

func main() {
	fmt.Println("[Client] start")

	config.Init()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("failed to start client: ", err)
		return
	}

	// id := 0
	for {
		// msg := []byte("hi, this is gozinx...")
		// msg := []byte(fmt.Sprintf("hi %v, this is gozinx", id))
		// id++

		// _, err := conn.Write(msg)
		// if err != nil {
		// 	fmt.Println("[Client] failed to write msg: ", err)
		// 	return
		// }

		// buf := make([]byte, 1024)
		// _, err = conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("[Client] failed to read msg: ", err)
		// 	return
		// }
		// fmt.Println("[Client] receive msg from server: ", string(buf))

		// 封包
		dp := datapack.New()

		data := []byte("hello, this is gozinx...")
		msg := &model.Message{
			MsgID:   0,
			DataLen: uint32(len(data)),
			Data:    data,
		}

		b, err := dp.Pack(msg)
		if err != nil {
			fmt.Printf("[Client] failed to pack msg: %v\n", err)
			return
		}
		if _, err := conn.Write(b); err != nil {
			fmt.Printf("[Client] failed to write msg: %v\n", err)
			return
		}

		// 服务端数据读取，有回复消息
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("[Client] failed to readfull head: ", err)
			break
		}
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("[Client] failed to UnPack head: ", err)
			break
		}

		if msgHead.DataLen > 0 {
			// 说明里面是有数据的
			msgHead.Data = make([]byte, msgHead.DataLen)
			if _, err := io.ReadFull(conn, msgHead.Data); err != nil {
				fmt.Println("[Client] failed to readfull msg data: ", err)
				return
			}
			fmt.Printf("[Client] recevie msg | msgID: %v, dataLen: %v, data: %v\n", msgHead.MsgID, msgHead.DataLen, string(msgHead.Data))
		}

		// test 休眠
		time.Sleep(1 * time.Second)
	}
}
