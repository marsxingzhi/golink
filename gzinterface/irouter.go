package gzinterface

type IRouter interface {
	PreHandle(req IRequest)
	Handle(re IRequest)
	PostHandle(req IRequest)
}
