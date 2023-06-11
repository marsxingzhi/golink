package router

import (
	"github.com/marsxingzhi/xzlink/pkg/request"
)

type IRouter interface {
	PreHandle(req request.IRequest)
	Handle(re request.IRequest)
	PostHandle(req request.IRequest)
}
