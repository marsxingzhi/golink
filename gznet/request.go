package gznet

import (
	"github.com/marsxingzhi/xzlink/model"
	conn "github.com/marsxingzhi/xzlink/pkg/connection"
)

// 链接与数据的封装
type Request struct {
	Conn conn.IConnection
	// Data []byte
	// 将Data封装到Message中
	Msg *model.Message
}

func (req *Request) GetConnection() conn.IConnection {
	return req.Conn
}

func (req *Request) GetData() []byte {
	return req.Msg.Data
}

func (req *Request) GetMsgID() uint32 {
	return req.Msg.MsgID
}
