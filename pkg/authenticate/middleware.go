package authenticate

import (
	"context"
	"github.com/gin-gonic/gin"
)

// GinMiddleware stores the logged-in user for further usage
func (a *Auth) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	Grab Authorization Cookie
		token := c.GetHeader("Authorization")
		if token != "" {
			//	Retrieve the related account to the token
			account, err := a.GetAccountByToken(context.Background(), token)
			if err != nil {
				c.Set("username", account.Username)
			}
		}
		c.Next()
	}
}
