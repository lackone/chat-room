package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/chat-room/helpers"
	"net/http"
)

// 判断用户是否登录
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "token不能为空",
			})
			return
		}
		jwtClaims, err := helpers.ParseJwt(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "token解析失败",
			})
			return
		}
		c.Set("jwt_claims", jwtClaims)
		c.Next()
	}
}
