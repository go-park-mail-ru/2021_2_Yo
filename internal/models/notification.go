package models

type Notification struct {
	Type        string
	ReceiverId  string
	UserId      string
	UserName    string
	UserSurname string
	UserImgUrl  string
	EventId     string
	EventTitle  string
}
