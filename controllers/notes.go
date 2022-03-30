package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNotes(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		return
	}
	fmt.Println("email:", email)
	filter := bson.M{"email": email}
	res := <-database.Api.FindDocs("notes", filter, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		fmt.Println("3Error:", res.Error.Error())
		return
	}
	cur := res.Result
	notes := make([]models.Note, 0)
	switch v := cur.(type) {
	case *mongo.Cursor:
		for v.Next(context.Background()) {
			note := &models.Note{}
			err := v.Decode(note)
			if err != nil {
				ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
				fmt.Println("4Error:", res.Error.Error())
				return
			}
			notes = append(notes, *note)
		}
	default:
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		fmt.Println("5Error:", res.Error.Error())
		return
	}

	if len(notes) <= 0 {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.NoNotes))
		fmt.Println("6Error:", res.Error.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(notes, constants.NoError))
	return
}

func AddNote(ctx *gin.Context) {
	type reqData struct {
		Data  string `json="data"`
		Title string `json="title"`
		Email string `json="email"`
	}
	rD := &reqData{}
	if err := ctx.BindJSON(rD); err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println(">>Error:", err.Error())
		return
	}
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println(">>Error:", errors.New("Email doesn't exist"))
		return
	}
	rD.Email = fmt.Sprintf("%v", email)
	mp := map[string]string{"data": rD.Data, "title": rD.Title, "email": rD.Email}
	note := models.Note{}
	n, err := models.Note.CreateNoteFromInterface(note, mp)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	res := <-database.Api.InsertDoc("notes", n)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	var id string
	switch v := res.Result.(type) {
	case gin.H:
		id = fmt.Sprintf("%v", v["_id"])
	default:
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		fmt.Println("5Error:", res.Error.Error())
		return
	}
	noteData := models.NoteData{}
	mp = map[string]string{"data": rD.Data, "parent": id, "email": rD.Email}
	nD, err := models.NoteData.CreateNoteDataViaMap(noteData, mp)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	ndRes := <-database.Api.InsertDoc("noteData", nD)
	if ndRes.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed))
		fmt.Println("Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(res.Result, constants.NoError))
}

func GetNoteData(ctx *gin.Context) {
	parent := ctx.Request.URL.Query().Get("parent")
	// if !exists {
	// 	ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
	// 	fmt.Println(">>Error:", errors.New("No parent ID found"))
	// 	return
	// }
	fmt.Println("Parent: ", parent)
	filter := bson.D{{"parent", parent}}
	noteData := &models.NoteData{}
	err := <-database.Api.FindDoc("noteData", filter, noteData, nil)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		fmt.Println(">>Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(noteData, constants.NoError))

}

func DeleteNotes(ctx *gin.Context) {

}

func RenameNote(ctx *gin.Context) {

}
