package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtToken struct {
	Token string
}
type User struct {
	FirstName string     `json:"firstName" binding:"required"`
	LastName  string     `json:"lastName" binding:"required"`
	Email     string     `json:"email" binding:"required,email" isUnique:"true"`
	Password  string     `json:"password"`
	IsActive  bool       `json:"isActive"`
	Tokens    []JwtToken `json:"tokens"`
}

func (u User) CreateUser(ctx *gin.Context) (*User, error) {
	user := &u
	if err := ctx.BindJSON(user); err != nil {
		return nil, err
	}
	user.IsActive = true
	return user, nil
}

func (u User) GenerateAuthToken() (*User, *string, error) {
	user := &u
	secretKeyStr := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	if len(secretKeyStr) <= 0 {
		err := errors.New("No JWT key found")
		return nil, nil, err
	}
	secretKey := []byte(secretKeyStr)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Key:", err)
		return nil, nil, err
	}
	jT := JwtToken{Token: tokenString}
	user.Tokens = append(user.Tokens, jT)
	return user, &jT.Token, nil
}
