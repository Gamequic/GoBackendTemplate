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
	r.POST("/signup", controllers.Signup) //perfil 4
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/users", middleware.RequireAuth, controllers.GetUsersList)      //perfil 3
	r.DELETE("/users/:id", middleware.RequireAuth, controllers.DeleteUser) //perfil 6

	r.Run() // listen and serve on 0.0.0.0:8080
}
