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
	r.POST("/api/users/signup", controllers.Signup)
	r.POST("/api/users/login", controllers.Login)
	r.GET("/api/users/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/api/users", middleware.RequireAuth, controllers.CreateUser)                        //perfil 4
	r.GET("/api/users", middleware.RequireAuth, controllers.GetUsers)                           //perfil 3
	r.GET("/api/users/:id", middleware.RequireAuth, controllers.GetUserById)                    //perfil 5
	r.PUT("/api/users/:id", middleware.RequireAuth, controllers.UpdateUser)                     //perfil 5
	r.DELETE("/api/users/:id", middleware.RequireAuth, controllers.DeleteUser)                  //perfil 6
	r.DELETE("/api/users/delete/:id", middleware.RequireAuth, controllers.DeleteUserCompletely) //perfil 1

	r.Run() // listen and serve on 0.0.0.0:8080
}
