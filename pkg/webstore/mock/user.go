package mock

import "github.com/sir-wiggles/arc/pkg/webstore"

// UserService is the mock service for user
type UserService struct {
	UserFn      func(id int) (*webstore.User, error)
	UserInvoked bool

	UsersFn      func() ([]*webstore.User, error)
	UsersInvoked bool

	CreateUserFn      func(u *webstore.User) error
	CreateUserInvoked bool

	DeleteUserFn      func(id int) error
	DeleteUserInvoked bool
}

// User mock
func (u *UserService) User(id int) (*webstore.User, error) {
	u.UserInvoked = true
	return u.UserFn(id)
}

// Users mock
func (u *UserService) Users() ([]*webstore.User, error) {
	u.UsersInvoked = true
	return u.UsersFn()
}

// CreateUser mock
func (u *UserService) CreateUser(user *webstore.User) error {
	u.CreateUserInvoked = true
	return u.CreateUserFn(user)
}

// DeleteUser mock
func (u *UserService) DeleteUser(id int) error {
	u.DeleteUserInvoked = true
	return u.DeleteUserFn(id)
}
