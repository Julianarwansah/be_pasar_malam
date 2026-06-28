package middleware

import (
	"strings"

	"pasarmalam/services"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "user_id"
	CtxEmail  = "user_email"
	CtxRole   = "user_role"
)

func Auth(jwtSvc *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Missing or invalid Authorization header",
				"error_code": "UNAUTHORIZED",
			})
			return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		claims, err := jwtSvc.Parse(tok)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"error_code": "INVALID_TOKEN",
			})
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
