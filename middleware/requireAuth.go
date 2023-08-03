package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mateenqazi/jwt-authenication/initializers"
	"github.com/mateenqazi/jwt-authenication/models"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("I am middleware")

	// Get the cookie off req

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/Validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECERT_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to req
		c.Set("user", user)

		// continue
		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
