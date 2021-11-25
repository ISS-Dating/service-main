package web

import (
	"github.com/ISS-Dating/service-main/model"
	"github.com/golang-jwt/jwt"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResgisterInfo struct {
}

type userClaims struct {
	model.User
	*jwt.StandardClaims
	TokenType string
}
