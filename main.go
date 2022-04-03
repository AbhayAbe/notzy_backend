package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AbhayAbe/notzy_backend/controllers"
	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var _PORT_ string

func main() {
	configEnv()
	database.ConfigureMongodb()
	handleConnectionAndRoutes()
}

func configEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Couldn't load env")
		os.Exit(1)
	}
	_PORT_ = ":" + os.Getenv("PORT")
}

func handleConnectionAndRoutes() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CorsWebsiteURL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowCredentials: true,
	}))

	aR := router.Group("/")

	//unauthenticated routes
	router.GET("/test", controllers.Test)

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	//authenticatedRoutes
	aR.Use(middlewares.Auth())
	{
		aR.GET("/authenticate", controllers.AuthenticateUser)
		aR.GET("/getUser", controllers.GetUser)
		aR.GET("/logout", controllers.Logout)
		aR.GET("/notes", controllers.GetNotes)
		aR.GET("/getNoteData", controllers.GetNoteData)

		aR.POST("/addNote", controllers.AddNote)
		aR.PATCH("/updateNote", controllers.UpdateNote)

		aR.PATCH("/saveNoteData", controllers.SaveNoteData)

		aR.DELETE("/note", controllers.DeleteNote)
		aR.DELETE("/user", controllers.DeleteUser)
	}

	if err := router.Run(_PORT_); err != nil {
		log.Fatal("Server couldn't start due to error:", err)
	}
	fmt.Println("Listening to server on port: ", _PORT_)
}
