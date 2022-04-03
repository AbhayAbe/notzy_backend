package utils

import "github.com/gin-gonic/gin"

func GenerateResponse(message interface{}, err string, additional interface{}) gin.H {
	return gin.H{
		"message":    message,
		"error":      err,
		"additional": additional,
	}

}
