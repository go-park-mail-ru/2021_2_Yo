package notificator

import (
	"backend/internal/models"
	"backend/internal/service/event"
	"backend/internal/service/notification"
	"backend/internal/service/notification/delivery/websocket"
	"backend/internal/service/user"
	log "github.com/sirupsen/logrus"
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
	Seen        bool   `json:"seen"`
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	UserSurname string `json:"userSurname"`
	UserImgUrl  string `json:"userImgUrl,omitempty"`
	EventId     string `json:"eventId,omitempty"`
	EventTitle  string `json:"eventTitle,omitempty"`
}

func (n *Notificator) NewSubscriberNotification(receiverId string, userId string) error {
	u, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	ws := n.pool.GetConn(receiverId)
	if ws != nil {
		m := &NotificationBody{
			Type:        "0",
			Seen:        false,
			UserId:      u.ID,
			UserName:    u.Name,
			UserSurname: u.Surname,
		}
		if u.ImgUrl != "" {
			m.UserImgUrl = u.ImgUrl
		}
		err := ws.WriteJSON(m)
		if err != nil {
			log.Error("NewSubscriberNotification err = ", err)
			n.pool.RemoveConn(receiverId)
			err = n.nRepository.CreateSubscribeNotification(receiverId, u, false)
			return err
		} else {
			err = n.nRepository.CreateSubscribeNotification(receiverId, u, false)
			return err
		}
	} else {
		log.Error("NewSubscriberNotification ws = ", ws)
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
	ws := n.pool.GetConn(receiverId)
	if ws != nil {
		m := &NotificationBody{
			Type:        "1",
			Seen:        false,
			UserId:      u.ID,
			UserName:    u.Name,
			UserSurname: u.Surname,
			EventId:     e.ID,
			EventTitle:  e.Title,
		}
		if u.ImgUrl != "" {
			m.UserImgUrl = u.ImgUrl
		}
		err := ws.WriteJSON(m)
		if err != nil {
			log.Error("InvitationNotification err = ", err)
			n.pool.RemoveConn(receiverId)
			err = n.nRepository.CreateInviteNotification(receiverId, u, e, false)
			return err
		} else {
			err = n.nRepository.CreateInviteNotification(receiverId, u, e, false)
			return err
		}
	} else {
		log.Error("InvitationNotification ws = ", ws)
		err = n.nRepository.CreateInviteNotification(receiverId, u, e, false)
		return err
	}
}

func (n *Notificator) NewEventNotification(userId string, eventId string) error {
	author, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	subscribers, err := n.uRepository.GetSubscribers(userId)
	if err != nil {
		return err
	}
	e, err := n.eRepository.GetEventById(eventId)
	if err != nil {
		return err
	}
	for _, sub := range subscribers {
		ws := n.pool.GetConn(sub.ID)
		if ws != nil {
			m := &NotificationBody{
				Type:        "2",
				Seen:        false,
				UserId:      author.ID,
				UserName:    author.Name,
				UserSurname: author.Surname,
				EventId:     e.ID,
				EventTitle:  e.Title,
			}
			if author.ImgUrl != "" {
				m.UserImgUrl = author.ImgUrl
			}
			err := ws.WriteJSON(m)
			if err != nil {
				n.pool.RemoveConn(sub.ID)
				err = n.nRepository.CreateNewEventNotification(sub.ID, author, e, false)
				return err
			} else {
				err = n.nRepository.CreateNewEventNotification(sub.ID, author, e, false)
				return err
			}
		} else {
			err = n.nRepository.CreateNewEventNotification(sub.ID, author, e, false)
			return err
		}
	}
	return nil
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
