package gznet

import (
	"github.com/marsxingzhi/golink/gzinterface"
	"github.com/marsxingzhi/golink/model"
)

// 链接与数据的封装
type Request struct {
	Conn gzinterface.IConnection
	// Data []byte
	// 将Data封装到Message中
	Msg *model.Message
}

func (req *Request) GetConnection() gzinterface.IConnection {
	return req.Conn
}

func (req *Request) GetData() []byte {
	return req.Msg.Data
}

func (req *Request) GetMsgID() uint32 {
	return req.Msg.MsgID
}
