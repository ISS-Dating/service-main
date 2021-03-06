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

func (s *Service) UpdateUser(user model.User) (model.User, error) {
	return s.Repo.UpdateUser(user)
}

func (s *Service) GetUserByUsername(username string) (model.User, error) {
	return s.Repo.ReadUserByUsername(username)
}

func (s *Service) ModUser(author model.User, username string, value bool) error {
	if author.Role != model.RoleAdministrator {
		return ErrorForbidden
	}

	user, err := s.Repo.ReadUserByUsername(username)
	if err != nil {
		return err
	}

	if value {
		user.Role = model.RoleModerator
	} else {
		user.Role = model.RoleUser
	}

	_, err = s.Repo.UpdateUser(user)

	return err
}

func (s *Service) BanUser(author model.User, username string, value bool) error {
	if author.Role != model.RoleAdministrator && author.Role != model.RoleModerator {
		return ErrorForbidden
	}

	user, err := s.Repo.ReadUserByUsername(username)
	if err != nil {
		return err
	}

	user.Banned = value
	user.Stats.BannedBefore = true

	_, err = s.Repo.UpdateUser(user)

	return err
}

func (s *Service) MatchUsers(usernameA, usernameB string) error {
	userA, _ := s.Repo.ReadUserByUsername(usernameA)
	userA.Stats.UsersMet++
	s.Repo.UpdateUser(userA)

	userB, _ := s.Repo.ReadUserByUsername(usernameB)
	userB.Stats.UsersMet++
	s.Repo.UpdateUser(userB)

	return s.Repo.CreateAcquaintance(usernameA, usernameB)
}

func (s *Service) ListFriends(username string) ([]string, error) {
	var list []string
	friends, err := s.Repo.ReadAcquaintanceByUsername(username)
	if err != nil {
		return nil, err
	}

	for _, f := range friends {
		if f.UserAUsername == username {
			list = append(list, f.UserBUsername)
		} else {
			list = append(list, f.UserAUsername)
		}
	}

	return list, nil
}
