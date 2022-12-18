package domain

type UserRepository interface {
	RegisterUser(user User, password string) error
}
