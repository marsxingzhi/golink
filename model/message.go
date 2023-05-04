package model

// 定义消息结构
type Message struct {
	// 消息id
	MsgID uint32
	// 消息长度
	DataLen uint32
	// 消息内容
	Data []byte
}
