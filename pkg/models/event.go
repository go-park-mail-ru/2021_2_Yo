package models

type Event struct {
	ID          string
	Title       string
	Description string
	Text        string
	City        string
	Category    string
	Viewed      int
	ImgUrl      string
	Tag         []string
	Date        string
	Geo         string
	Address     string
	AuthorId    string
}
