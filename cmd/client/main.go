package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("start client")

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("failed to start client: ", err)
		return
	}

	for {
		msg := []byte("hi, this is gozinx...")
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("[Client] failed to write msg: ", err)
			return
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("[Client] failed to read msg: ", err)
			return
		}
		fmt.Println("[Client] receive msg from server: ", string(buf))

		// test 休眠
		time.Sleep(1 * time.Second)
	}
}
