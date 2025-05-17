package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// миддлваре убеждается что в запросе содержится валидный JWT и сохраняет его данные в контекст
func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Auth header!"})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("role", claims.Role)
		c.Next()
	}
}

// миддлваре чтобы проверить что роль пользователя совпадает с требуемой ролью(ролями)
func Authorize(allowedRoles ...string) gin.HandlerFunc {
	// для удобства перевожу список в мапу
	roleSet := map[string]struct{}{}
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		v, exists := c.Get("role")
		if !exists { // нет требований для доступа к запросу <==> кто угодно может воспользоваться
			c.Next()
			return
		}
		role, _ := v.(string)

		if _, allowed := roleSet[role]; !allowed {
			// переданной в токене роли нету в списке тех кому можно разрешить использовать этот endpoint <==> forbidden
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён"})
			return
		}

		c.Next()
	}
}
