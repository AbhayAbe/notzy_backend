package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	docId := ctx.Request.URL.Query().Get("id")
	limit := ctx.Request.URL.Query().Get("limit")
	page := ctx.Request.URL.Query().Get("page")
	sort := ctx.Request.URL.Query().Get("sort")
	var pageLimit int64 = 0
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
		return
	}
	fmt.Println("email:", email)
	filter := bson.M{"email": email}
	if len(docId) > 0 {
		id, err := primitive.ObjectIDFromHex(docId)
		if err != nil {
			ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
			println("%%Error:", err.Error())
			return
		}

		filter["_id"] = id
	}
	opts := options.Find()
	opts.SetProjection(bson.M{"data": 0})
	if len(sort) > 0 {
		split := strings.Split(sort, ",")
		arr := make([][]string, 0)
		for _, s := range split {
			sp := strings.Split(s, ":")
			arr = append(arr, sp)
		}
		sortPattern := bson.M{}
		for _, ele := range arr {
			i, err := strconv.ParseInt(ele[1], 10, 64)
			if err != nil {
				ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
				println("%%Error:", err.Error())
				return
			}
			sortPattern[ele[0]] = i
		}

		opts.SetSort(sortPattern)
	}
	if len(limit) > 0 {
		i, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
			println("%%Error:", err.Error())
			return
		}
		pageLimit = i
		opts.SetLimit(i)
	}

	if len(page) > 0 {
		i, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
			println("%%Error:", err.Error())
			return
		}
		fmt.Println("Currpage:", i, ",limit:", pageLimit)
		if i > 1 && pageLimit > 0 {
			opts.SetSkip(i + pageLimit)
		}
	}
	res := <-database.Api.FindDocs("notes", filter, opts)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
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
				ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
				fmt.Println("4Error:", res.Error.Error())
				return
			}
			notes = append(notes, *note)
		}
	default:
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
		fmt.Println("5Error:", res.Error.Error())
		return
	}

	if len(notes) <= 0 {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.NoNotes, nil))
		fmt.Println("6Error:", res.Error.Error())
		return
	}
	docCount, err := database.DB.Collection("notes").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.NoNotes, nil))
		fmt.Println("6Error:", res.Error.Error())
		return
	}
	fmt.Println("DocCount:", docCount)
	ctx.JSON(200, utils.GenerateResponse(notes, constants.NoError, bson.M{"maxDocs": docCount}))
	return
}

func AddNote(ctx *gin.Context) {
	note := models.Note{}
	n, err := models.Note.CreateNote(note, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed, nil))
		fmt.Println("Error:", err.Error())
		return
	}
	n.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	n.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	res := <-database.Api.InsertDoc("notes", n)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.SaveNoteFailed, nil))
		fmt.Println("Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(res.Result, constants.NoError, nil))
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
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("!##Error:", err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("%%Error:", err.Error())
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"data": data.Data, "updatedAt": primitive.Timestamp{T: uint32(time.Now().Unix())}}}

	res := <-database.Api.UpdateDoc("notes", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("%%Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(constants.UpdateSuccesful, constants.NoError, nil))
}

func UpdateNote(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	fmt.Println("_ID:", docId)
	type noteStruct struct {
		Title     string              `json:"title" binding:"required"`
		UpdatedAt primitive.Timestamp `json:"updatedAt"`
	}
	note := &noteStruct{}

	err := ctx.BindJSON(note)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("!##Error:", err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("%%Error:", err.Error())
		return
	}
	filter := bson.M{"_id": id}
	fmt.Println("Filter:", filter)
	note.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	bM, err := bson.Marshal(note)
	bD := &bson.D{}
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("!##Error:", err.Error())
		return
	}
	err = bson.Unmarshal(bM, &bD)
	update := bson.M{
		"$set": bD,
	}
	res := <-database.Api.UpdateDoc("notes", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.UpdateNotefailed, nil))
		println("!##Error:", res.Error.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(note, constants.NoError, nil))
}

func GetNoteData(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
		println("%%Error:", err.Error())
		return
	}
	filter := bson.M{"_id": id}
	noteData := &models.Note{}
	err = <-database.Api.FindDoc("notes", filter, noteData, nil)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.GetNotesFailed, nil))
		fmt.Println(">>Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(noteData, constants.NoError, nil))
}

func DeleteNote(ctx *gin.Context) {
	docId := ctx.Request.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", err.Error())
		return
	}
	noteFilter := bson.M{"_id": id}
	res := <-database.Api.DeleteDoc("notes", noteFilter, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", err.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(constants.DeleteSuccesful, constants.NoError, nil))
}
