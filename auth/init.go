package auth

import "errors"

var (
	SecretKey []byte
	ExpTime   int64
)

const (
	errTokenNotValid  = "ERROR_TOKEN_NOT_VALID"
	errSecretIsWeak   = "ERROR_SECRET_IS_WEAK"
	errUnAuth         = "UN_AUTHORIZED"
	errNoPerms        = "NO_PERMISSIONS"
	errTokenIsExpired = "TOKEN_EXPIRED"
)

const (
	UserRoleAdmin        = "ADMIN"
	InternalServicesRole = "INTERNAL_SERVICES"
	SuperUserRole        = "SUPER_USER"
)

//Init manual
func Init(secret string, expTime int64) error {

	// check secret length
	if len(secret) < 10 {
		return errors.New(errSecretIsWeak)
	}

	// set secret key
	SecretKey = []byte(secret)

	// set expire time
	ExpTime = expTime
	return nil
}
