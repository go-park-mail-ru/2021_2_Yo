package csrf

type Manager interface {
	Create(userId string) (string, error)
	Check(sessionId string) (string, error)
	Delete(sessionId string) error
}
