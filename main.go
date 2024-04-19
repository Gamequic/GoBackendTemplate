package main

import (
	"github.com/calleros/go-jwt/controllers"
	"github.com/calleros/go-jwt/initializers"
	"github.com/calleros/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.LoadDbTables(initializers.DB)

}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/users", controllers.GetUsersList)

	r.Run() // listen and serve on 0.0.0.0:8080
}
