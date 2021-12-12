package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

const (
	PUBLISH   = "publish"
	SUBSCRIBE = "subscribe"
)

type Pool struct {
	mutex       sync.RWMutex
	Connections map[string]*websocket.Conn
}

func NewPubSub() *Pool {
	return &Pool{
		Connections: make(map[string]*websocket.Conn),
	}
}

func (p *Pool) AddConn(userId string, ws *websocket.Conn) {
	p.mutex.Lock()
	p.Connections[userId] = ws
	p.mutex.Unlock()
}

func (p *Pool) RemoveConn(userId string) {
	p.mutex.Lock()
	p.Connections[userId] = nil
	p.mutex.Unlock()
}

func (p *Pool) GetConn(userId string) *websocket.Conn {
	return p.Connections[userId]
}
