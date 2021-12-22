package websocket

import (
	"backend/internal/response"
	log "backend/pkg/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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
	_, ok := w.(http.Hijacker)
	if !ok {
		log.Info(!ok)
	}
	userId := r.Context().Value(response.CtxString("userId")).(string)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	p.AddConn(userId, conn)
	log.Info("WebsocketHandler new client with id: ", userId, " total clients: ", len(p.Connections))
}
