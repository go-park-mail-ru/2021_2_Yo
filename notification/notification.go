package notification

import (
	"backend/easyWebsocket"
	"strconv"
)

type SubsNotificator struct {
	userConnectionsPool *easyWebsocket.PubSub
}

func NewSubsNotificator (userPool *easyWebsocket.PubSub) *SubsNotificator {
	return &SubsNotificator{
		userConnectionsPool: userPool,
	}
}

func (sn * SubsNotificator) NewSubscriber(subId string, subbedToId string) {
	_, err := strconv.Atoi(subId)
	if err != nil {
		return
	}
}