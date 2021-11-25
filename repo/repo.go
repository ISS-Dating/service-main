package repo

import (
	"errors"

	"github.com/ISS-Dating/service-main/model"
)

var (
	ErrorUsernameDuplication = errors.New("username is already taken")
	ErrorUserNotExist        = errors.New("user does not exist")
)

type Interface interface {
	CreateUser(user model.User) (model.User, error)
	ReadUserByLogin(username, password string) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
	// ReadUserByUsername(username string) (*model.User, error)

	// CreateAcquaintance(userA, userB string) error
	// GetAcquaintanceByUsername(username string) ([]model.Acquaintance, error)
}
