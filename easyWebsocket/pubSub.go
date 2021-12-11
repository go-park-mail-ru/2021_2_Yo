package easyWebsocket

import (
	"net/http"
	"strconv"
	"sync"
	log "backend/pkg/logger"
	"github.com/gorilla/websocket"
)

type PubSub struct {
	mutex sync.RWMutex
	Connections map[int]*websocket.Conn
}

func NewPubSub() *PubSub {
	return &PubSub{
		Connections: make (map[int]*websocket.Conn),
	}
}

func (p *PubSub) AddConn(userId int, ws *websocket.Conn) {
	p.mutex.Lock()
	p.Connections[userId] = ws
	p.mutex.Unlock()
}

func (p *PubSub) RemoveConn (userId int) {
	p.mutex.Lock()
	p.Connections[userId] = nil
	p.mutex.Unlock()
}

func (p *PubSub) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		return 
	}

	userId, err := GetID(ws)
	if err != nil {
		return
	}

	iUserId, _ := strconv.Atoi(userId)

	p.AddConn(iUserId,ws)
	
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Error("Something went wrong", err)

			p.RemoveConn(iUserId)
			//store notifications

			return
		}
		//need to do write funcs
		ws.HandleReceiveMessage(client, messageType, p)
	}
}