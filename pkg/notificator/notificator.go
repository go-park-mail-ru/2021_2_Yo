package notificator

import (
	"backend/internal/models"
	"backend/internal/service/event"
	"backend/internal/service/notification"
	"backend/internal/service/notification/delivery/websocket"
	"backend/internal/service/user"
	"time"
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

func (n *Notificator) createAndSendNotification(notification *NotificationBody, receiverId string, user *models.User, event *models.Event, repoFunc func(string, *models.User, *models.Event) error) error {
	ws := n.pool.GetConn(receiverId)
	if ws != nil {
		err := ws.WriteJSON(notification)
		if err != nil {
			n.pool.RemoveConn(receiverId)
			err = repoFunc(receiverId, user, event)
			return err
		} else {
			err = repoFunc(receiverId, user, event)
			return err
		}
	} else {
		err := repoFunc(receiverId, user, event)
		return err
	}
}

func (n *Notificator) NewSubscriberNotification(receiverId string, userId string) error {
	u, err := n.uRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	nf := &NotificationBody{
		Type:        "0",
		Seen:        false,
		UserId:      u.ID,
		UserName:    u.Name,
		UserSurname: u.Surname,
	}
	if u.ImgUrl != "" {
		nf.UserImgUrl = u.ImgUrl
	}
	return n.createAndSendNotification(nf, receiverId, u, nil, n.nRepository.CreateSubscribeNotification)
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
	return n.createAndSendNotification(m, receiverId, u, e, n.nRepository.CreateInviteNotification)
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
	for _, sub := range subscribers {
		err := n.createAndSendNotification(m, sub.ID, author, e, n.nRepository.CreateNewEventNotification)
		if err != nil {
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

func (n *Notificator) EventTomorrowNotification() error {
	currentTime := time.Now()
	currentDate := currentTime.Format("02.01.2006")
	events, err := n.eRepository.GetEvents("", "", "", "", currentDate, nil)
	if err != nil {
		return err
	}
	for _, e := range events {
		visitors, err := n.uRepository.GetVisitors(e.ID)
		if err != nil {
			return err
		}
		author, err := n.uRepository.GetUserById(e.AuthorId)
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
		for _, v := range visitors {
			err := n.createAndSendNotification(m, v.ID, author, e, n.nRepository.CreateTomorrowEventNotification)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
