package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//UserClaims struct
type UserClaims struct {
	jwt.StandardClaims
	UserID   interface{} `json:"uid"`
	Username *string     `json:"username,omitempty"`
	AppID    *string     `json:"app_id,omitempty"`
	Email    *string     `json:"email,omitempty"`
	Status   *string     `json:"status,omitempty"`
	Roles    []string    `json:"roles"`
}

//GenerateToken generate jwt token
func (u *UserClaims) GenerateToken() (*string, error) {

	// set expire time
	u.ExpiresAt = time.Now().Unix() + expTime

	// new with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

//CheckRoles check access roles
func (u *UserClaims) CheckRoles(roles ...string) bool {

	// for-loop check []vs[] token roles
	for _, r := range roles {
		for _, tr := range u.Roles {
			if r == tr || tr == UserRoleAdmin {
				return true
			}
		}
	}

	return false
}
