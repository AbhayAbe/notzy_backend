package controllers

import (
	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database/mongo/api/notzyMongo/crud"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user models.User
	u, err := models.User.CreateUser(user, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println(err)
		return
	}
	res := <-crud.InsertDoc("users", u)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println(err)
		return
	}
	ctx.JSON(200, utils.GenerateResponse(u, constants.NoError))
}

func Login(ctx *gin.Context) {

}

func Logout(ctx *gin.Context) {

}

func GoogleSignin(ctx *gin.Context) {

}

func AppleSignin(ctx *gin.Context) {

}
