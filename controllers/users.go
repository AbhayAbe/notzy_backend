package controllers

import (
	"github.com/gin-gonic/gin"
)

func DeleteUser(ctx *gin.Context) {
	// email, exists := ctx.Get("email")
	// if !exists{
	// 	ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed))
	// 	println("%%Error:", err.Error())
	// 	return
	// }
	// filter := bson.D{{"email", docId}}
	// res := <-database.Api.DeleteDoc("noteData", noteDataFilter, nil)
	// if res.Error != nil {
	// 	ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed))
	// 	println("%%Error:", err.Error())
	// 	return
	// }
	// res = <-database.Api.DeleteDoc("notes", noteFilter, nil)
	// if res.Error != nil {
	// 	ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed))
	// 	println("%%Error:", err.Error())
	// 	return
	// }
	// ctx.JSON(200, utils.GenerateResponse(constants.DeleteSuccesful, constants.NoError))
}

func UpdateUser(ctx *gin.Context) {

}

func DeactivateUser(ctx *gin.Context) {

}

func ActivateUser(ctx *gin.Context) {

}

func DeleteUsers(ctx *gin.Context) {

}
