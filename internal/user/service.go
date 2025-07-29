package user

type Service interface {
	Register(user *User, password string) error
	Authenticate(email, password string) (*User, error)
	Authorize(user *User, role string) bool
}
