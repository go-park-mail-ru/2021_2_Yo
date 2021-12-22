package response

import (
	models "backend/internal/models"
)

const logMessage = "response:response:"

type HttpStatus int
type CtxString string

type Response struct {
	Status  HttpStatus  `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

type UserResponseBody struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty" valid:"type(string),length(0|50)" san:"xss"`
	Surname  string `json:"surname,omitempty" valid:"type(string),length(0|50)" san:"xss"`
	About    string `json:"description,omitempty" valid:"type(string),length(0|150)" san:"xss"`
	ImgUrl   string `json:"imgUrl,omitempty" valid:"type(string)" san:"xss"`
	Mail     string `json:"email,omitempty" valid:"email,length(0|150)" san:"xss"`
	Password string `json:"password,omitempty" valid:"type(string),length(0|150)" san:"xss"`
}

type UserListResponseBody struct {
	Users []UserResponseBody `json:"users"`
}

type EventIDResponseBody struct {
	ID string `json:"id"`
}

type EventResponseBody struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title" valid:"type(string),length(0|520)" san:"xss"`
	Description string   `json:"description" valid:"type(string),length(0|1020)" san:"xss"`
	Text        string   `json:"text" valid:"type(string),length(0|5000)" san:"xss"`
	City        string   `json:"city" valid:"type(string),length(0|60)" san:"xss"`
	Category    string   `json:"category" valid:"type(string),length(0|30)" san:"xss"`
	Viewed      int      `json:"viewed" valid:"type(int)" san:"xss"`
	ImgUrl      string   `json:"imgUrl" valid:"type(string),length(0|255)" san:"xss"`
	Tag         []string `json:"tag" san:"xss"`
	Date        string   `json:"date" valid:"type(string),length(0|10)" san:"xss"`
	Geo         string   `json:"geo" valid:"type(string),length(0|255)"`
	Address     string   `json:"address" valid:"type(string), length(0|520)" san:"xss"`
	AuthorID    string   `json:"authorid" san:"xss"`
	IsVisited   bool     `json:"favourite"`
}

type EventListResponseBody struct {
	Events []EventResponseBody `json:"events"`
}

type SubscribedResponseBody struct {
	Result bool `json:"result"`
}

type FavouriteResponseBody struct {
	Result bool `json:"result"`
}

type CitiesResponseBody struct {
	Cities []string `json:"cities"`
}

type NotificationResponseBody struct {
	Type        string `json:"type"`
	Seen        bool   `json:"seen"`
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	UserSurname string `json:"userSurname"`
	UserImgUrl  string `json:"userImgUrl"`
	EventId     string `json:"eventId,omitempty"`
	EventTitle  string `json:"eventTitle,omitempty"`
}

type NotificationListResponseBody struct {
	Notifications []NotificationResponseBody `json:"notifications"`
}

func StatusResponse(status HttpStatus) *Response {
	return &Response{
		Status: status,
	}
}

func OkResponse() *Response {
	return &Response{
		Status:  200,
		Message: "OK",
	}
}

func UserResponse(user *models.User) *Response {
	return &Response{
		Status: 200,
		Body:   MakeUserResponseBody(user),
	}
}

func UserListResponse(users []*models.User) *Response {
	return &Response{
		Status: 200,
		Body:   MakeUserListResponseBody(users),
	}
}

func EventResponse(event *models.Event) *Response {
	return &Response{
		Status: 200,
		Body:   MakeEventResponseBody(event),
	}
}

func EventIdResponse(eventID string) *Response {
	return &Response{
		Status: 200,
		Body: EventIDResponseBody{
			ID: eventID,
		},
	}
}

func EventListResponse(events []*models.Event) *Response {
	return &Response{
		Status: 200,
		Body:   MakeEventListResponseBody(events),
	}
}

func SubscribedResponse(result bool) *Response {
	return &Response{
		Status: 200,
		Body: SubscribedResponseBody{
			Result: result,
		},
	}
}

func FavouriteResponse(result bool) *Response {
	return &Response{
		Status: 200,
		Body: FavouriteResponseBody{
			Result: result,
		},
	}
}

func CitiesResponse(cities []string) *Response {
	return &Response{
		Status: 200,
		Body: CitiesResponseBody{
			Cities: cities,
		},
	}
}

func NotificationListResponse(notifications []*models.Notification) *Response {
	return &Response{
		Status: 200,
		Body:   MakeNotificationListResponseBody(notifications),
	}
}
