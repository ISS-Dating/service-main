package service

import (
	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/repo"
)

type Service struct {
	Repo repo.Interface
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(username, password string) (model.User, error) {
	return model.User{}, nil
}

func (s *Service) Register(username, password, email string) (model.User, error) {
	return s.Repo.CreateUser(model.User{
		Username: username,
		Password: password,
		Email:    email,
	})
}
