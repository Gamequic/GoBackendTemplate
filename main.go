package main

import (
	"github.com/calleros/sich/controllers"
	"github.com/calleros/sich/initializers"
	"github.com/calleros/sich/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.LoadUsersTables(initializers.DB)

}

func main() {
	r := gin.Default()
	r.POST("/api/users/signup", controllers.Signup)
	r.POST("/api/users/login", controllers.Login)
	r.POST("/api/users/logout", controllers.Logout)
	r.GET("/api/users/validate", controllers.Validate)
	r.POST("/api/users", middleware.RequireAuth, controllers.CreateUser)
	r.GET("/api/users", middleware.RequireAuth, controllers.GetUsers)
	r.GET("/api/users/:id", middleware.RequireAuth, controllers.GetUserById)
	r.PUT("/api/users/:id", middleware.RequireAuth, controllers.UpdateUser)
	r.DELETE("/api/users/:id", middleware.RequireAuth, controllers.DeleteUser)
	r.DELETE("/api/users/:id/delete", middleware.RequireAuth, controllers.DeleteUserCompletely)
	r.POST("/api/profiles", middleware.RequireAuth, controllers.CreateProfile)
	r.GET("/api/profiles", middleware.RequireAuth, controllers.GetProfiles)
	r.GET("/api/profiles/:id", middleware.RequireAuth, controllers.GetProfileById)
	r.PUT("/api/profiles/:id", middleware.RequireAuth, controllers.UpdateProfile)
	r.GET("/api/profiles/:id/users", middleware.RequireAuth, controllers.GetUsersByProfileId)

	r.Run() // listen and serve on 0.0.0.0:8080
}
