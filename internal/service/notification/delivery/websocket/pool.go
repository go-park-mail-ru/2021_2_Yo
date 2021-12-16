package websocket

import (
	log "backend/pkg/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Pool struct {
	mutex       sync.RWMutex
	Connections map[string]*websocket.Conn
}

func NewPool() *Pool {
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

func (p *Pool) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("HERE1")

	userID, err := GetID(conn)
	if err != nil {
		log.Error(err)
	}
	log.Info("HERE2")

	p.AddConn(userID, conn)
	log.Info("New Client is connected with id: ", userID, "total: ", len(p.Connections))
}
