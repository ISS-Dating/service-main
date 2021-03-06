package web

import (
	"github.com/ISS-Dating/service-main/model"
	"github.com/golang-jwt/jwt"
)

type genericRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Mod      bool   `json:"mod"`
	Ban      bool   `json:"ban"`
}

type Status struct {
	Status bool `json:"status"`
}

type MatchedList struct {
	Friends []string `json:"friends"`
}

type userClaims struct {
	model.User
	*jwt.StandardClaims
	TokenType string
}
