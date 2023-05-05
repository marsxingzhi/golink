package handler

import (
	"fmt"

	"github.com/marsxingzhi/gozinx/gzinterface"
)

// 消息管理模块
type IMsghandler interface {
	// 处理消息
	DoHandle(request gzinterface.IRequest)
	// 添加路由
	AddRouter(msgID uint32, router gzinterface.IRouter)
}

type MsgHandler struct {
	Apis map[uint32]gzinterface.IRouter
}

func New() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]gzinterface.IRouter),
	}
}

func (h *MsgHandler) DoHandle(request gzinterface.IRequest) {
	router, ok := h.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("[MsgHandler] DoHandle | not found router with msgID: %v\n", request.GetMsgID())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (h *MsgHandler) AddRouter(msgID uint32, router gzinterface.IRouter) {
	_, ok := h.Apis[msgID]
	if ok {
		fmt.Printf("[MsgHandler] AddRouter | has be existed request with msgID: %v\n", msgID)
		return
	}
	h.Apis[msgID] = router
}
