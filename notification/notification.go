package notification

import (
	"backend/websocket"
)

type Notificator struct {
	pool *websocket.Pool
}

func NewSubsNotificator(pool *websocket.Pool) *Notificator {
	return &Notificator{
		pool: pool,
	}
}

func (sn *Notificator) NewSubscriber(subscriberId string, subscribedName string) error {
	ws := sn.pool.GetConn(subscriberId)
	type Message struct {
		Name string `json:"Name"`
	}
	m := &Message{
		Name: subscribedName,
	}
	err := ws.WriteJSON(m)
	if err != nil {
		sn.pool.RemoveConn(subscriberId)
		return err
	}
	return nil
}
