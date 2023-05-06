package gznet

import (
	"fmt"
	"sync"

	"github.com/marsxingzhi/golink/gzinterface"
)

type ConnectionManager struct {
	conns map[uint32]gzinterface.IConnection
	mutex sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns: make(map[uint32]gzinterface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnectionManager) Add(conn gzinterface.IConnection) {
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	connMgr.conns[conn.GetConnID()] = conn

	fmt.Printf("[ConnectionManager] Add connID = :%v successfully, and now connection total num: %v\n", conn.GetConnID(), connMgr.Len())
}

// 删除链接
func (connMgr *ConnectionManager) Remove(conn gzinterface.IConnection) {
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	delete(connMgr.conns, conn.GetConnID())
	fmt.Printf("[ConnectionManager] Remove connID: %v successfully, and now connection total num: %v\n", conn.GetConnID(), connMgr.Len())
}

// 获取链接
func (connMgr *ConnectionManager) Get(connID uint32) (gzinterface.IConnection, bool) {
	// 读锁
	connMgr.mutex.RLock()
	defer connMgr.mutex.RUnlock()

	conn, ok := connMgr.conns[connID]
	return conn, ok
}

// 链接个数
func (connMgr *ConnectionManager) Len() int {
	return len(connMgr.conns)
}

// 清理所有链接
func (connMgr *ConnectionManager) ClearAllConns() {
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	for connID, conn := range connMgr.conns {

		// 1. 关闭资源
		// 2. 删除
		conn.Stop()

		delete(connMgr.conns, connID) // 能在这里直接删除吗？ 边遍历，边删除？
	}
	fmt.Println("[ConnectionManager] has cleared all connections...")
}
