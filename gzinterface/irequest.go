package gzinterface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}