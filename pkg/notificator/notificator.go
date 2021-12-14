package notificator

import (
	"backend/internal/models"
	"backend/internal/service/event"
	"backend/internal/service/notification"
	"backend/internal/service/notification/delivery/websocket"
	"backend/internal/service/user"
)

type Notificator struct {
	pool        *websocket.Pool
	nRepository notification.Repository
	uRepository user.Repository
	eRepository event.Repository
}

func NewNotificator(pool *websocket.Pool, nr notification.Repository, ur user.Repository, er event.Repository) *Notificator {
	return &Notificator{
		pool:        pool,
		nRepository: nr,
		uRepository: ur,
		eRepository: er,
	}
}

type NotificationBody struct {
	Type        string `json:"type"`
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	UserSurname string `json:"userSurname"`
	UserImgUrl  string `json:"userImgUrl"`
	EventId     string `json:"eventId,omitempty"`
	EventTitle  string `json:"eventTitle,omitempty"`
}

func (n *Notificator) NewSubscriberNotification(receiverId string, userId string) error {
	u, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	ws := n.pool.GetConn(userId)
	if ws != nil {
		m := &NotificationBody{
			Type:        "subscription",
			UserId:      u.ID,
			UserName:    u.Name,
			UserSurname: u.Surname,
			UserImgUrl:  u.ImgUrl,
		}
		err := ws.WriteJSON(m)
		if err != nil {
			n.pool.RemoveConn(userId)
			err = n.nRepository.CreateSubscribeNotification(receiverId, u, false)
			return err
		} else {
			err = n.nRepository.CreateSubscribeNotification(receiverId, u, true)
			return err
		}
	} else {
		err := n.nRepository.CreateSubscribeNotification(receiverId, u, false)
		return err
	}
}

func (n *Notificator) DeleteSubscribeNotification(receiverId string, userId string) error {
	return n.nRepository.DeleteSubscribeNotification(receiverId, userId)
}

func (n *Notificator) InvitationNotification(receiverId string, userId string, eventId string) error {
	u, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	e, err := n.eRepository.GetEventById(eventId)
	if err != nil {
		return err
	}
	ws := n.pool.GetConn(userId)
	if ws != nil {
		m := &NotificationBody{
			Type:        "subscription",
			UserId:      u.ID,
			UserName:    u.Name,
			UserSurname: u.Surname,
			UserImgUrl:  u.ImgUrl,
			EventId:     e.ID,
			EventTitle:  e.Title,
		}
		err := ws.WriteJSON(m)
		if err != nil {
			n.pool.RemoveConn(userId)
			err = n.nRepository.CreateInviteNotification(receiverId, u, e, false)
			return err
		} else {
			err = n.nRepository.CreateInviteNotification(receiverId, u, e, true)
			return err
		}
	} else {
		err = n.nRepository.CreateInviteNotification(receiverId, u, e, false)
		return err
	}
}

func (n *Notificator) NewEventNotification(receiverId string, userId string, eventId string) error {
	u, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	e, err := n.eRepository.GetEventById(eventId)
	if err != nil {
		return err
	}
	ws := n.pool.GetConn(userId)
	if ws != nil {
		m := &NotificationBody{
			Type:        "subscription",
			UserId:      u.ID,
			UserName:    u.Name,
			UserSurname: u.Surname,
			UserImgUrl:  u.ImgUrl,
			EventId:     e.ID,
			EventTitle:  e.Title,
		}
		err := ws.WriteJSON(m)
		if err != nil {
			n.pool.RemoveConn(userId)
			err = n.nRepository.CreateNewEventNotification(receiverId, u, e, false)
			return err
		} else {
			err = n.nRepository.CreateNewEventNotification(receiverId, u, e, true)
			return err
		}
	} else {
		err = n.nRepository.CreateNewEventNotification(receiverId, u, e, false)
		return err
	}
}

func (n *Notificator) UpdateNotificationsStatus(receiverId string) error {
	return n.nRepository.UpdateNotificationsStatus(receiverId)
}

func (n *Notificator) GetAllNotifications(receiverId string) ([]*models.Notification, error) {
	return n.nRepository.GetAllNotifications(receiverId)
}

func (n *Notificator) GetNewNotifications(receiverId string) ([]*models.Notification, error) {
	return n.nRepository.GetNewNotifications(receiverId)
}