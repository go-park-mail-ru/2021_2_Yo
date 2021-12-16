package postgres

import "backend/internal/models"

type Notification struct {
	Type        string `db:"type"`
	ReceiverId  string `db:"receiver_id"`
	UserId      string `db:"user_id"`
	UserName    string `db:"user_name"`
	UserSurname string `db:"user_surname"`
	UserImgUrl  string `db:"user_img_url"`
	EventId     string `db:"event_id"`
	EventTitle  string `db:"event_title"`
	Seen        bool   `db:"seen"`
}

func toModelNotification(n *Notification) *models.Notification {
	return &models.Notification{
		Type:        n.Type,
		ReceiverId:  n.ReceiverId,
		UserId:      n.UserId,
		UserName:    n.UserName,
		UserSurname: n.UserSurname,
		UserImgUrl:  n.UserImgUrl,
		EventId:     n.EventId,
		EventTitle:  n.EventTitle,
	}
}
