package service

import (
	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/repo"
)

type Service struct {
	Repo repo.Interface
}

func NewService(repo repo.Interface) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Login(username, password string) (model.User, error) {
	return s.Repo.ReadUserByLogin(username, password)
}

func (s *Service) Register(username, password, email string) (model.User, error) {
	return s.Repo.CreateUser(model.User{
		Username: username,
		Password: password,
		Email:    email,
	})
}
