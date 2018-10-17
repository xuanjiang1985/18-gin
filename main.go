package main

import (
	"18-gin/api/user"
	"18-gin/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(utils.Cors())
	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", user.Name)
	router.Run(":8091")
}
