package controllers

import (
	"fmt"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database/mongo/api/notzyMongo/crud"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
)

func GetNotes(ctx *gin.Context) {

}

func AddNote(ctx *gin.Context) {
	note := models.Note{}
	fmt.Println(ctx.Get("data"))
	n, err := models.Note.CreateNote(note, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	res := <-crud.InsertDoc("notes", n)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(res.Result, constants.NoError))
}

func DeleteNotes(ctx *gin.Context) {

}

func RenameNote(ctx *gin.Context) {

}
