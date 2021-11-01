package models

type ResponseBodyUser struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty" valid:"type(string)" san:"xss"`
	Surname  string `json:"surname,omitempty" valid:"type(string)" san:"xss"`
	About    string `json:"description,omitempty" valid:"type(string)" san:"xss"`
	Mail     string `json:"email,omitempty" valid:"email" san:"xss"`
	Password string `json:"password,omitempty" valid:"type(string)" san:"xss"`
}

type ResponseBodyEvent struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title" san:"xss"`
	Description string   `json:"description" san:"xss"`
	Text        string   `json:"text" san:"xss"`
	City        string   `json:"city" san:"xss"`
	Category    string   `json:"category" san:"xss"`
	Viewed      int      `json:"viewed" san:"xss"`
	ImgUrl      string   `json:"imgUrl" san:"xss"`
	Tag         []string `json:"tag" san:"xss"`
	Date        string   `json:"date" san:"xss"`
	Geo         string   `json:"geo" san:"xss"`
	AuthorID    string   `json:"authorid" san:"xss"`
}
