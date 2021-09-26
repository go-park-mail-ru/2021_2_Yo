package auth

type UseCase interface {
	SignUp(name, surname, mail,password string) error
	SignIn(name, surname, mail,password string) (string, error)
	Parse(cookie string) (string, error)
	List() []string
}
