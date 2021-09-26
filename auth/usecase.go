package auth

type UseCase interface {
	SignUp(username, password string) error
	SignIn(username, password string) (string, error)
	//ParseToken(accessToken string) (*models.User, error)
	List() []string
}
