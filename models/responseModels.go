package models

type ResponseBodyUser struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty" valid:"type(string)"`
	Surname  string `json:"surname,omitempty" valid:"type(string)"`
	About    string `json:"description,omitempty" valid:"type(string)"`
	Mail     string `json:"email,omitempty" valid:"email"`
	Password string `json:"password,omitempty" valid:"type(string)"`
}

type ResponseBodyEvent struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Text        string   `json:"text"`
	City        string   `json:"city"`
	Category    string   `json:"category"`
	Viewed      int      `json:"viewed"`
	ImgUrl      string   `json:"imgUrl"`
	Tag         []string `json:"tag"`
	Date        string   `json:"date"`
	Geo         string   `json:"geo"`
	AuthorID    string   `json:"authorID"`
}
