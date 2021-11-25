package service

import "github.com/ISS-Dating/service-main/model"

type Interface interface {
	Login(username, password string) (model.User, error)
	Register(username, password, email string) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
}
