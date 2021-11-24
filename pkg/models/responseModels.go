package models

type UserResponseBody struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty" valid:"type(string),length(0|50)" san:"xss"`
	Surname  string `json:"surname,omitempty" valid:"type(string),length(0|50)" san:"xss"`
	About    string `json:"description,omitempty" valid:"type(string),length(0|150)" san:"xss"`
	ImgUrl   string `json:"imgUrl,omitempty" valid:"type(string)" san:"xss"`
	Mail     string `json:"email,omitempty" valid:"email,length(0|150)" san:"xss"`
	Password string `json:"password,omitempty" valid:"type(string),length(0|50)" san:"xss"`
}

type UserListResponseBody struct {
	Users []UserResponseBody `json:"users"`
}

type EventIDResponseBody struct {
	ID string `json:"id"`
}

type EventResponseBody struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title" valid:"type(string),length(0|255)" san:"xss"`
	Description string   `json:"description" valid:"type(string),length(0|500)" san:"xss"`
	Text        string   `json:"text" valid:"type(string),length(0|2200)" san:"xss"`
	City        string   `json:"city" valid:"type(string),length(0|30)" san:"xss"`
	Category    string   `json:"category" valid:"type(string),length(0|30)" san:"xss"`
	Viewed      int      `json:"viewed" valid:"type(int)" san:"xss"`
	ImgUrl      string   `json:"imgUrl" valid:"type(string),length(0|255)" san:"xss"`
	Tag         []string `json:"tag" san:"xss"`
	Date        string   `json:"date" valid:"type(string),length(0|10)" san:"xss"`
	Geo         string   `json:"geo" valid:"type(string),length(0|255)"`
	Address		string	 `json:"address" valid:"type(string)","length(0|255)" san:"xss"` 
	AuthorID    string   `json:"authorid" san:"xss"`
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
