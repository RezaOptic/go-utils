package auth

import "errors"

var (
	secretKey []byte
	// TODO: need refresh token for reNew jwt
	expTime int64 = 8640000 // 100 days
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
func Init(secret string) error {

	// check secret length
	if len(secret) < 10 {
		return errors.New(errSecretIsWeak)
	}

	// set secret key
	secretKey = []byte(secret)

	return nil
}
