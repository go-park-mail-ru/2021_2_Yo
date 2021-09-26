package auth

type UseCase interface {
	SignUp(username, password string) error
	SignIn(username, password string) (string, error)
	Parse(cookie string) (string, error)
	List() []string
}
