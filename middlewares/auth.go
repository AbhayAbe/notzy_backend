package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Next()
		authToken := strings.Replace(ctx.Request.Header.Get("Authorization"), "Bearer ", "", -1)

		secretKeyStr := os.Getenv("JWT_SECRET")
		if len(secretKeyStr) <= 0 {
			ctx.JSON(403, constants.AuthenticationFailed)
			ctx.Abort()
		}
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Validation failed")
			}
			return []byte(secretKeyStr), nil
		})
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)["email"]
			ctx.Set("email", claims)
			ctx.Next()
		} else {
			fmt.Println("Error:", err)
			ctx.JSON(403, constants.AuthenticationFailed)
			ctx.Abort()
		}
	}
}
