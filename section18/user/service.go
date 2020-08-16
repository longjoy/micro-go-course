package user

import "context"

type UserService interface {
	CheckPassword(ctx context.Context, username string, password string) (bool, error)
}

type UserServiceImpl struct{}

func (s UserServiceImpl) CheckPassword(ctx context.Context, username string, password string) (bool, error) {
	if username == "admin" && password == "admin" {
		return true, nil
	}
	return false, nil
}
