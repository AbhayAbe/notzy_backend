package models

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email" isUnique:"true"`
	Password  string `json:"password" binding:"required"`
	IsActive  bool   `json:"isActive"`
}

func (u User) CreateUser(ctx *gin.Context) (*User, error) {
	user := &User{}
	if err := ctx.BindJSON(user); err != nil {
		return nil, err
	}
	user.IsActive = true
	return user, nil
}
