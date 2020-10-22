package jwt

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

//必须具有角色role
func RequiredRole(role int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		account := c.MustGet(util.TokenKey).(*util.Claims)

		if account.Role != role {
			app.UnauthorizedResp(c, e.ERROR_AUTH, fmt.Sprintf("required role[%d]", role))
			c.Abort()
			return
		}
	}
}

// 多个角色中只要有一个匹配就算通过
func RequiredManyOfOneRole(roles ...int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		account := c.MustGet(util.TokenKey).(*util.Claims)

		flag := false
		for _, role := range roles {
			if account.Role == role {
				flag = true
				break
			}
		}

		if !flag {
			app.UnauthorizedResp(c, e.ERROR_AUTH, fmt.Sprintf("required role[%d]", roles))
			c.Abort()
			return
		}
	}
}
