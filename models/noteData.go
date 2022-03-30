package models

import (
	"github.com/gin-gonic/gin"
)

type NoteData struct {
	Id     interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	Parent string      `json:"parent" binding:"required" isUnique:"true"`
	Data   string      `json:"data" binding:"required"`
}

func (d NoteData) CreateNoteData(ctx *gin.Context) (*NoteData, error) {
	data := &d
	if err := ctx.BindJSON(data); err != nil {
		return nil, err
	}
	return data, nil
}
func (d NoteData) CreateNoteDataViaMap(note map[string]string) (*NoteData, error) {
	data := &d
	data.Parent = note["parent"]
	data.Data = note["data"]
	return data, nil
}
