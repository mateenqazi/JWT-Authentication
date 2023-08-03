package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenqazi/jwt-authenication/controllers"
	"github.com/mateenqazi/jwt-authenication/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	r.Run()
}
