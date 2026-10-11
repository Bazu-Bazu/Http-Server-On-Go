package user

type Repository interface {
	Create(u User) (User, error)
	GetByID(id int64) (User, error)
	Delete(id int64) error
}
