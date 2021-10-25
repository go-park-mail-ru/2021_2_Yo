package models

type Event struct {
	ID          string
	Title       string
	Description string
	Text        string
	City        string
	Category    string
	Viewed      int
	Tag         []string
	Date        string
	Geo         string
}
