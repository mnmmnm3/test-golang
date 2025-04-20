package middleware

import (
	"fmt"
	"go-api/initializers"
	"go-api/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	// get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		return
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return
		}

		// find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			return
		}

		// attach to req
		c.Set("user", user)

		// continue
		c.Next()
	} else {
		return
	}
}
