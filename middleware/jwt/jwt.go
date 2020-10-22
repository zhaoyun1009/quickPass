package jwt

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/token_cache"
	"QuickPass/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID          int64  `json:"userId"`
	AccountName string `json:"accountName"`
	RoleId      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	jwt.StandardClaims
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(util.HeaderToken)
		if token == "" {
			app.UnauthorizedResp(c, e.ERROR_AUTH, "")
			c.Abort()
			return
		}

		claims, err := util.ParseToken(token)
		if err != nil {
			code := e.ERROR_AUTH_CHECK_TOKEN_FAIL
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
			app.UnauthorizedResp(c, code, "")
			c.Abort()
			return
		}

		// 判断jwt是否还有5分钟过期，如果快过期，则签发新的token到响应头
		// jwt过期时间
		expireTime := claims.ExpiresAt - time.Now().Unix()
		if expireTime < setting.AppSetting.TokenRenewalTime {
			cacheToken, err := token_cache.GetCacheToken(claims.Agency, claims.Username)
			if err == nil {
				// token相同，表明当前token没有续签,需要续签
				if cacheToken == "" || cacheToken == token {
					cacheToken, _ = util.GenerateToken(claims.Agency, claims.Username, claims.Nickname, claims.Role)
				}
				c.Writer.Header().Add(util.HeaderToken, cacheToken)
			}
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set(util.TokenKey, claims)
	}
}
