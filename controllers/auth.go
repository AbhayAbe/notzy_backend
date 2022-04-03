package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AbhayAbe/notzy_backend/constants"
	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/statics"
	"github.com/AbhayAbe/notzy_backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthenticateUser(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.AuthenticationFailed, nil))
		return
	}
	fmt.Println("email:", email)
	ctx.JSON(200, utils.GenerateResponse(constants.Authenticationsuccesful, "", nil))
	return
}
func Register(ctx *gin.Context) {
	var user models.User
	var token *string
	u, err := models.User.CreateUser(user, ctx)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed, nil))
		println("!##Error:", err.Error())
		return
	}
	u, token, err = models.User.GenerateAuthToken(*u)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed, nil))
		println("!##Error:", err.Error())
		return
	}
	hash, err := utils.GenerateFromPassword(u.Password, statics.ArgonParams)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed, nil))
		println("##Error:", err)
		return
	}
	u.Password = hash
	res := <-database.Api.InsertDoc("users", u)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed, nil))
		println("##Error:", err)
		return
	}
	u.Password = ""
	u.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(gin.H{"user": u, "token": *token}, constants.NoError, nil))
}

func Login(ctx *gin.Context) {
	data := make(map[string]string)
	jsonData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		fmt.Println("Error:", err)
		return
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		fmt.Println("Error:", err)
		return
	}
	user := &models.User{}
	email := data["email"]
	filter := bson.D{{"email", email}}
	er := <-database.Api.FindDoc("users", filter, user, nil)
	if er != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		fmt.Println("Error:", err)
		return
	}

	match, err := utils.ComparePasswordAndHash(data["password"], user.Password)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		fmt.Println("Error:", err)
		return
	}
	if !match {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Match: %v\n", match)

	var token *string
	user, token, err = models.User.GenerateAuthToken(*user)

	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.RegistrationFailed, nil))
		println("!##Error:", err.Error())
		return
	}
	update := bson.D{
		{"$set", bson.D{{"tokens", user.Tokens}}},
	}
	res := <-database.Api.UpdateDoc("users", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LoginFailed, nil))
		println("##Error:", res.Error.Error())
		return
	}
	user.Password = ""
	user.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(gin.H{"user": user, "token": *token}, constants.NoError, nil))
	fmt.Println("Error:", err)
	return
}

func Logout(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed, nil))
		return
	}
	fmt.Println("email:", email)
	user := &models.User{}
	filter := bson.D{{"email", email}}
	err := <-database.Api.FindDoc("users", filter, user, nil)
	if err != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed, nil))
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
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed, nil))
		fmt.Println("Couldn't remove tokens")
		return
	}
	user.Tokens = fTokens
	update := bson.D{
		{"$set", bson.D{{"tokens", fTokens}}},
	}
	res := <-database.Api.UpdateDoc("users", filter, update, nil)
	if res.Error != nil {
		ctx.JSON(500, utils.GenerateResponse(nil, constants.LogoutFailed, nil))
		println("##Error:", res.Error.Error())
		return
	}
	user.Password = ""
	user.Tokens = make([]models.JwtToken, 0)
	ctx.JSON(200, utils.GenerateResponse(user, "", nil))
}

func GoogleSignin(ctx *gin.Context) {

}

func AppleSignin(ctx *gin.Context) {

}
