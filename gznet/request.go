package gznet

import "github.com/marsxingzhi/gozinx/gzinterface"

// 链接与数据的封装
type Request struct {
	Conn gzinterface.IConnection
	Data []byte
}

func (req *Request) GetConnection() gzinterface.IConnection {
	return req.Conn
}

func (req *Request) GetData() []byte {
	return req.Data
}
