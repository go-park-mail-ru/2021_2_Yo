package notification

import (
	"backend/easyWebsocket"
)

type SubsNotificator struct {
	userConnectionsPool *easyWebsocket.PubSub
}

func NewSubsNotificator(userPool *easyWebsocket.PubSub) *SubsNotificator {
	return &SubsNotificator{
		userConnectionsPool: userPool,
	}
}

func (sn *SubsNotificator) NewSubscriber(subberId string, subberName string) error {
	ws := sn.userConnectionsPool.GetConn(subberId)
	type Message struct {
		Name string `json:"Name`
	}
	m := &Message{
		Name: subberName,
	}
	err := ws.WriteJSON(m)
	if err != nil {
		sn.userConnectionsPool.RemoveConn(subberId)
		return err
	}
	return nil
}
