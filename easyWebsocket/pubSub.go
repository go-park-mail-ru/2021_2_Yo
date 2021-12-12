package easyWebsocket

import (
	//log "backend/pkg/logger"
	"encoding/json"
	//"net/http"
	//"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	PUBLISH   = "publish"
	SUBSCRIBE = "subscribe"
)

type Message struct {
	Action string `json:"action"`
	//Для подписки нужен id на кого подписались
	SubbedTo string `json:"subbedTo"`
	//На запас
	Message json.RawMessage `json:"message"`
}

type PubSub struct {
	mutex       sync.RWMutex
	Connections map[string]*websocket.Conn
}

func NewPubSub() *PubSub {
	return &PubSub{
		Connections: make(map[string]*websocket.Conn),
	}
}

func (p *PubSub) AddConn(userId string, ws *websocket.Conn) {
	p.mutex.Lock()
	p.Connections[userId] = ws
	p.mutex.Unlock()
}

func (p *PubSub) RemoveConn(userId string) {
	p.mutex.Lock()
	p.Connections[userId] = nil
	p.mutex.Unlock()
}

func (p *PubSub) GetConn(userId string) *websocket.Conn {
	return p.Connections[userId]
}

/*
func (p *PubSub) Subscribe(userId string, subbedTo string) {
	Просто отправь данные лол
}

/*
func(p *PubSub) HandleRecieveMessage(userId string, messageType int,  message []byte) (*PubSub) {
	m := Message{}

	err := json.Unmarshal(message, &m)
	if err != nil {
		log.Error("not correct json")
		return p
	}

	switch m.Action {

	case SUBSCRIBE:
		p.Subscribe(userId,m.SubbedTo)

	}
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
		//Вот это вообще не нужно по идее
		//p.HandleRecieveMessage(userId, messageType, message)
	}
}
*/
