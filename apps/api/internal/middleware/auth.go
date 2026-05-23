package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	pkgauth "github.com/zhaozeguang/timecapsule-api/pkg/auth"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{Code: 40100, Message: "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.Response{Code: 40100, Message: "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := pkgauth.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			code := 40100
			if err == pkgauth.ErrTokenExpired {
				code = 40101
			}
			c.JSON(http.StatusUnauthorized, dto.Response{Code: code, Message: err.Error()})
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, dto.Response{Code: 40100, Message: "invalid token type"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
