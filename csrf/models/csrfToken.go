package models

type CSRFData struct {
	CSRFToken  string
	UserId     string
	Expiration int
}
