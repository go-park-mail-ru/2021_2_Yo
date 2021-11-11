package jwtCsrf

type Manager interface {
	Create(userId string) (string, error)
	Check(susToken string) (string, error)
}
