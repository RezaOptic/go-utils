# go-utils

This is utils for Go project.

## auth:
a package for generate and validate access token.
access token contains the following structure.
```go
UserID   interface{}
Username *string    
AppID    *string    
Email    *string    
Status   *string    
Roles    []string   
``` 

### Usage:

---

### Initialize:

Import it in your code:

```go
import "github.com/RezaOptic/go-utils/auth"
```

at the first, we need to initialize the package for configs.
```go
auth.Init(ACCESS_TOKEN_SECRET, EXPIRE_TIME_IN_SECOND)
```

### Generate token:
```go
userClaim := auth.UserClaims{UserID: UserID}
Token, err := userClaim.Token()
if err != nil {
	panic(err)
}
fmt.Println(Token)
```

### Validate token:

#### **ginMiddleware**:

 for adding protection to your route add gin Middleware like this:
 ```go
 Foo.GET("/foo/user/:user_id", auth.JwtAuth(BarController, SELF_ACCESS, OPTIONAL, ROLES))
 ```
 and in your `BarController` can get user token information like this:
 ```go
 UserTokenInfotmation := auth.GetClaimFromContext(c)
 ```

 `SELF_ACCESS`: if you have `:user_id` param in your path for authorization TokenUserID and `user_id` must be the same.

 `OPTIONAL`: if you set this flag to true sending `Authorization` in the header is optional and if you don't send the token `auth.GetClaimFromContext(c)` return an empty struct.

 `ROLES`: roles is an array of string token must include these roles.
