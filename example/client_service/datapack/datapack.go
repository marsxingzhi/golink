package datapack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	clientConfig "github.com/marsxingzhi/xzlink/example/client_service/config"
	"github.com/marsxingzhi/xzlink/pkg/model"
)

// 封包、拆包逻辑
// 消息格式：[消息长度][消息编号][消息内容]
// 消息长度占4字节；消息编号占4字节

type DataPack struct {
}

func New() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// 消息长度 + 消息编号
	return 8
}

func (dp *DataPack) Pack(msg *model.Message) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// 1. 写入消息长度 写入二进制数据
	err := binary.Write(buffer, binary.LittleEndian, msg.DataLen)
	if err != nil {
		fmt.Printf("[DataPack] Pack | failed to write datalen: %v\n", err)
		return nil, err
	}
	// 2. 写入消息编号
	err = binary.Write(buffer, binary.LittleEndian, msg.MsgID)
	if err != nil {
		fmt.Printf("[DataPack] Pack | failed to write msgID: %v\n", err)
		return nil, err
	}
	// 3. 写入消息内容
	err = binary.Write(buffer, binary.LittleEndian, msg.Data)
	if err != nil {
		fmt.Printf("[DataPack] Pack | failed to write data: %v\n", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (dp *DataPack) UnPack(b []byte) (*model.Message, error) {
	// 1. 先读取head
	reader := bytes.NewReader(b)
	var msg model.Message
	// 这里不是读取全部，只是读取长度
	err := binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		fmt.Printf("[DataPack] UnPack | failed to read dataLen: %v\n", err)
		return nil, err
	}

	if err = binary.Read(reader, binary.LittleEndian, &msg.MsgID); err != nil {
		fmt.Printf("[DataPack] UnPack | failed to read msgID: %v\n", err)
		return nil, err
	}

	// 2. 判断dataLen的长度是否超过最大允许的包大小
	if clientConfig.Config.GetMaxPackageSize() > 0 && clientConfig.Config.GetMaxPackageSize() < int32(msg.DataLen) {
		fmt.Printf("[DataPack] UnPack | dataLen too large\n")
		return nil, errors.New("msg received is too large")
	}

	// 3. 然后根据head中的长度，读取消息内容。这个放到外面读
	// if err = binary.Read(reader, binary.LittleEndian, &msg.Data); err != nil {
	// 	fmt.Printf("[DataPack] UnPack | failed to read data: %v\n", err)
	// 	return nil, err
	// }
	return &msg, nil
}
