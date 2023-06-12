package datapack

import (
	"github.com/marsxingzhi/xzlink/pkg/model"
)

type IDataPack interface {
	// 获取包头长度
	GetHeadLen() uint32
	// 封包
	Pack(msg model.Message) ([]byte, error)
	// 拆包
	UnPack(b []byte) (model.Message, error)
}
