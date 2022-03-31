package controllers

import (
	"context"
	"fmt"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNotes(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		return
	}
	fmt.Println("email:", email)
	filter := bson.M{"email": email}
	opts := options.Find().SetProjection(bson.D{{"data", 0}})
	res := <-database.Api.FindDocs("notes", filter, opts)
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
	note := models.Note{}
	n, err := models.Note.CreateNote(note, ctx)
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
	ctx.JSON(200, utils.GenerateResponse(res.Result, constants.NoError))
}

func SaveNoteData(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	fmt.Println("_ID:", docId)

	type noteData struct {
		Data string `json:"data" binding:"required"`
	}
	data := &noteData{}
	err := ctx.BindJSON(data)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("!##Error:", err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("%%Error:", err.Error())
		return
	}

	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{{"data", data.Data}}},
	}

	res := <-database.Api.UpdateDoc("notes", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("%%Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(constants.UpdateSuccesful, constants.NoError))
}

func UpdateNote(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	fmt.Println("_ID:", docId)
	note := models.Note{}
	n, err := models.Note.CreateNote(note, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("!##Error:", err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("%%Error:", err.Error())
		return
	}
	filter := bson.M{"_id": id}
	fmt.Println("Filter:", filter)
	bM, err := bson.Marshal(n)
	bD := &bson.D{}
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("!##Error:", err.Error())
		return
	}
	err = bson.Unmarshal(bM, &bD)
	update := bson.D{
		{"$set", bD},
	}
	res := <-database.Api.UpdateDoc("notes", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed))
		println("!##Error:", res.Error.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(n, constants.NoError))
}

func GetNoteData(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		println("%%Error:", err.Error())
		return
	}
	filter := bson.D{{"_id", id}}
	noteData := &models.Note{}
	err = <-database.Api.FindDoc("notes", filter, noteData, nil)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed))
		fmt.Println(">>Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(noteData, constants.NoError))
}

func DeleteNote(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed))
		println("%%Error:", err.Error())
		return
	}
	noteFilter := bson.M{"_id": id}
	res := <-database.Api.DeleteDoc("notes", noteFilter, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed))
		println("%%Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(constants.DeleteSuccesful, constants.NoError))
}
