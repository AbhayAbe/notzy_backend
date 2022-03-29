package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database/mongo/api/notzyMongo/crud"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/statics"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(ctx *gin.Context) {
	var user models.User
	var token *string
	u, err := models.User.CreateUser(user, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println("!##Error:", err.Error())
		return
	}
	u, token, err = models.User.GenerateAuthToken(*u)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println("!##Error:", err.Error())
		return
	}
	hash, err := utils.GenerateFromPassword(u.Password, statics.ArgonParams)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println("##Error:", err)
		return
	}
	u.Password = hash
	res := <-crud.InsertDoc("users", u)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println("##Error:", err)
		return
	}
	u.Password = ""
	u.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(gin.H{"user": u, "token": *token}, constants.NoError))
}

func Login(ctx *gin.Context) {
	data := make(map[string]string)
	jsonData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		fmt.Println("Error:", err)
		return
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		fmt.Println("Error:", err)
		return
	}
	user := &models.User{}
	email := data["email"]
	filter := bson.D{{"email", email}}
	er := <-crud.FindDoc("users", filter, user, nil)
	if er != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		fmt.Println("Error:", err)
		return
	}

	match, err := utils.ComparePasswordAndHash(data["password"], user.Password)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		fmt.Println("Error:", err)
		return
	}
	if !match {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Match: %v\n", match)

	var token *string
	user, token, err = models.User.GenerateAuthToken(*user)

	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed))
		println("!##Error:", err.Error())
		return
	}
	update := bson.D{
		{"$set", bson.D{{"tokens", user.Tokens}}},
	}
	res := <-crud.UpdateDoc("users", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed))
		println("##Error:", res.Error.Error())
		return
	}
	user.Password = ""
	user.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(gin.H{"user": user, "token": *token}, constants.NoError))
	fmt.Println("Error:", err)
	return
}

func Logout(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		return
	}
	fmt.Println("email:", email)
	user := &models.User{}
	filter := bson.D{{"email", email}}
	err := <-crud.FindDoc("users", filter, user, nil)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		fmt.Println("##Error:", err.Error())
		return
	}
	tokens := user.Tokens
	var fTokens []models.JwtToken = make([]models.JwtToken, 0)
	authToken := strings.Replace(ctx.Request.Header.Get("Authorization"), "Bearer ", "", -1)
	isRemoved := false
	for idx, t := range tokens {
		if t.Token == authToken {
			fmt.Println("Token found")
			fTokens = append(tokens[:idx], tokens[idx+1:]...)
			fmt.Println(idx)
			isRemoved = true
			break
		}
	}
	if !isRemoved {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		fmt.Println("Couldn't remove tokens")
		return
	}
	user.Tokens = fTokens
	update := bson.D{
		{"$set", bson.D{{"tokens", fTokens}}},
	}
	res := <-crud.UpdateDoc("users", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed))
		println("##Error:", res.Error.Error())
		return
	}
	user.Password = ""
	user.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(user, ""))
}

func GoogleSignin(ctx *gin.Context) {

}

func AppleSignin(ctx *gin.Context) {

}
