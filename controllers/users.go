package controllers

import (
	"fmt"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteUser(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	type pass struct {
		Password string `json:"password"`
	}
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", constants.DeleteFailed)
		return
	}
	password := &pass{}
	err := ctx.BindJSON(password)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", err.Error())
		return
	}

	user := &models.User{}
	filter := bson.D{{"email", email}}
	er := <-database.Api.FindDoc("users", filter, user, nil)
	if er != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		fmt.Println("Error:", er.Error())
		return
	}
	fmt.Println("Password:", password.Password)
	match, err := utils.ComparePasswordAndHash(password.Password, user.Password)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.AuthenticationFailed, nil))
		fmt.Println("Error:", err)
		return
	}
	if !match {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.AuthenticationFailed, nil))
		fmt.Println("Error:", err)
		return
	}

	res := <-database.Api.DeleteDocs("notes", filter, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", res.Error.Error())
		return
	}
	res = <-database.Api.DeleteDoc("users", filter, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		println("%%Error:", res.Error.Error())
		return
	}
	ctx.JSON(200, utils.GenerateResponse(constants.DeleteSuccesful, constants.NoError, nil))
}

func GetUser(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.AuthenticationFailed, nil))
		return
	}
	user := &models.User{}
	filter := bson.D{{"email", email}}
	er := <-database.Api.FindDoc("users", filter, user, nil)
	if er != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.DeleteFailed, nil))
		fmt.Println("Error:", er.Error())
		return
	}
	user.Password = ""
	user.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(user, constants.NoError, nil))
}

func DeactivateUser(ctx *gin.Context) {

}

func ActivateUser(ctx *gin.Context) {

}

func DeleteUsers(ctx *gin.Context) {

}
