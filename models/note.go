package models

import (
	"errors"
	"fmt"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
)

type Note struct {
	Id    interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string      `json:"title" binding:"required"`
	Email string      `json:"email"`
}

func (n Note) CreateNote(ctx *gin.Context) (*Note, error) {
	note := &n
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		return nil, errors.New("Email doesn't exist")
	}
	if err := ctx.BindJSON(note); err != nil {
		return nil, err
	}
	note.Email = fmt.Sprintf("%v", email)
	return note, nil
}

func (n Note) CreateNoteFromInterface(data map[string]string) (*Note, error) {
	note := &n
	note.Title = data["title"]
	note.Email = data["email"]
	return note, nil
}
