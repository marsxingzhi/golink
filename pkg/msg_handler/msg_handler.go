package msg_handler

import (
	"fmt"
	"github.com/marsxingzhi/xzlink/pkg/config"
	"github.com/marsxingzhi/xzlink/pkg/request"
	"github.com/marsxingzhi/xzlink/pkg/router"
)

// 消息管理模块
// 接口里面是对外暴露的方法
type IMsgHandler interface {
	// 处理消息
	DoHandle(request request.IRequest)
	// 添加路由
	AddRouter(msgID uint32, router router.IRouter)
	StartWorkerPool()
	SendMessageToTaskQueue(request request.IRequest)
}

type MsgHandler struct {
	Apis map[uint32]router.IRouter

	// 消息队列，一个消息队列对应一个Worker
	TaskQueue []chan request.IRequest
	// Worker工作池的大小，即里面有多少个Worker
	WorkerPoolSize int
}

func New() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]router.IRouter),
		TaskQueue:      make([]chan request.IRequest, config.Config.GetWorkerPoolSize()),
		WorkerPoolSize: config.Config.GetWorkerPoolSize(),
	}
}

func (h *MsgHandler) DoHandle(request request.IRequest) {
	router, ok := h.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("[MsgHandler] DoHandle | not found router with msgID: %v\n", request.GetMsgID())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (h *MsgHandler) AddRouter(msgID uint32, router router.IRouter) {
	_, ok := h.Apis[msgID]
	if ok {
		fmt.Printf("[MsgHandler] AddRouter | has be existed request with msgID: %v\n", msgID)
		return
	}
	h.Apis[msgID] = router
}

// 启动工作池
func (h *MsgHandler) StartWorkerPool() {
	// 1. 根据Size开启worker，并开启流程
	// 2. 给消息队列开启空间
	// 3. 启动工作流程

	for i := 0; i < h.WorkerPoolSize; i++ {
		// 肯定是有缓冲的channel
		h.TaskQueue[i] = make(chan request.IRequest, config.Config.GetWorkerTaskCapacity())

		go h.StartWorker(i, h.TaskQueue[i])
	}
}

// 启动worker工作流程
// 工作流程：
// 1. 阻塞等待从消息队列中获取消息，即request
// 2. 将request交给MsgHandler执行
func (h *MsgHandler) StartWorker(workerID int, taskQueue chan request.IRequest) {
	fmt.Println("[MsgHandler] StartWorker | workderID: ", workerID, " is starting...")
	for {
		select {
		case request := <-taskQueue:
			h.DoHandle(request)
		}
	}
}

// SendMessageToTaskQueue 将消息发送到消息队列，然后消息队列将消息交给Worker处理
// 流程：
// 1. 将消息平均分配给不同的Worker
// 1.1 根据requestID分配给Worker（由于不存在requestID，这里先使用ConnectionID，后续改进） TODO
// 2. 将消息发送给对应的Worker的TaskQueue即可
func (h *MsgHandler) SendMessageToTaskQueue(request request.IRequest) {
	// TODO 目前只是单体应用，还不是分布式的，如果是分布式的话，需要使用hash一致性算法
	workerID := request.GetConnection().GetConnID() % uint32(h.WorkerPoolSize)

	fmt.Println("[MshHandler] SendMessageToTaskQueue | workerID: ", workerID)

	h.TaskQueue[workerID] <- request

}
