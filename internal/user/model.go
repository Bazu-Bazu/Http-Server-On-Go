package user

type User struct {
	ID           int64
	Nickname     string
	Name         string
	PasswordHash string
}

type RegisterInput struct {
	Nickname string
	Name     string
	Password string
}
