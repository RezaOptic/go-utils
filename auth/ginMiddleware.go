package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	claimKey = "claims"
)


//GinJwtAuth gin middleware for authentication jwt token
func GinJwtAuth(function gin.HandlerFunc, selfAccess, optional bool, roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check jwt
		token := strings.Replace(ctx.GetHeader("Authorization"), "Bearer ", "", -1)
		if token == "" && optional {
			// token is valid
			ctx.Set(claimKey, &UserClaims{})
			// call
			function(ctx)
			ctx.Next()
		} else {
			claims, err := ValidateToken(token)
			// optional flag skip jwt authentication and set a empty UserClaims model to header
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": errUnAuth})
				return
			}

			// check access roles
			var skipSelfAccess bool
			if len(roles) > 0 {
				if ok := claims.CheckRoles(roles...); !ok && !selfAccess {
					ctx.JSON(http.StatusForbidden, gin.H{"error": errNoPerms})
					return
				} else if ok { // role is OK
					skipSelfAccess = true
				}
			}

			// check selfy role with `user_id` params
			// `user_id` locked params
			if selfAccess && !skipSelfAccess {
				uidParams := ctx.Param("user_id")
				if claims.UserID != uidParams {
					ctx.JSON(http.StatusForbidden, gin.H{"error": errNoPerms})
					return
				}
			}

			// token is valid
			ctx.Set(claimKey, *claims)

			// call
			function(ctx)
			ctx.Next()
		}
	}
}

//GetClaimFromContext get claim from gin context
func GetClaimFromContext(ctx *gin.Context) *UserClaims {
	resp, exist := ctx.Get(claimKey)
	if !exist {
		return nil
	}

	c := resp.(UserClaims)
	return &c
}
