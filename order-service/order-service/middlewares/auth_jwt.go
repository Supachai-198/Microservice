package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthJWT() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {

		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			defer c.AbortWithStatus(http.StatusUnauthorized)
		}

		tokenHeader := c.GetHeader("Authorization")
		accessToken := strings.SplitAfter(tokenHeader, "Bearer")[1]
		// fmt.Println(accessToken)
		jwtSecretKey := os.Getenv("JWT_KEY")

		token, _ := jwt.Parse(strings.Trim(accessToken, " "), func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			defer c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			// global value result
			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", claims["user_id"])
			// return to next method if token is exist
			c.Next()
		}

	})
}
