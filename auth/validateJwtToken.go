package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)


//ValidateToken validate token
func ValidateToken(tokenString string) (*UserClaims, error) {

	// parse with claims
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New(errTokenNotValid)
	}

	// cast to custom userClaims
	tokenClaims := token.Claims.(*UserClaims)

	// check expire time
	if tokenClaims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New(errTokenIsExpired)
	}

	return tokenClaims, nil
}
