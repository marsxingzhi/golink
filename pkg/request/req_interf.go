package request

import (
	"github.com/marsxingzhi/xzlink/pkg/connection"
)

type IRequest interface {
	GetConnection() connection.IConnection
	GetData() []byte
	GetMsgID() uint32
}
