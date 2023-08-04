package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenqazi/jwt-authenication/controllers"
	_ "github.com/mateenqazi/jwt-authenication/docs"
	"github.com/mateenqazi/jwt-authenication/initializers"
	"github.com/mateenqazi/jwt-authenication/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

// @title Tag Service API
// @version 1.0
// @description A Tag service API in Go using Gin framework

// @host localhost:3000
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validation", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
