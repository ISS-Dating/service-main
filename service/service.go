package service

import (
	"errors"

	"github.com/ISS-Dating/service-main/model"
)

var (
	ErrorForbidden = errors.New("Forbidden")
)

type Interface interface {
	Login(username, password string) (model.User, error)
	Register(username, password, email string) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
	GetUserByUsername(username string) (model.User, error)
	ModUser(author model.User, username string, value bool) error
	BanUser(author model.User, username string, value bool) error
	MatchUsers(usernameA, usernameB string) error
	ListFriends(username string) ([]string, error)
}
