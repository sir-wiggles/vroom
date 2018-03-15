package webstore

type User struct {
	ID        int
	FirstName string
	LastName  string
	Phone     int64
}

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id int) error
}
