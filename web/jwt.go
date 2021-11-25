package web

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"time"

	"github.com/ISS-Dating/service-main/model"
	"github.com/golang-jwt/jwt"
)

var (
	signKey     *rsa.PrivateKey
	validateKey *rsa.PublicKey
)

func readToken(tokenStr string) (model.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateKey, nil
	})
	if err != nil {
		return model.User{}, err
	}

	claims := token.Claims.(*userClaims)
	return claims.User, nil
}

func createToken(user model.User) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &userClaims{
		user,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
		"tokenize",
	}

	return t.SignedString(signKey)
}

func auth(req *http.Request) (model.User, int) {
	c, err := req.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return model.User{}, http.StatusUnauthorized
		}
		return model.User{}, http.StatusBadRequest
	}

	token := c.Value
	user, err := readToken(token)
	if err != nil {
		return model.User{}, http.StatusBadRequest
	}

	return user, http.StatusOK
}
