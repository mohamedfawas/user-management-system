package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-management-system/controllers"
	"github.com/mohamedfawas/user-management-system/handlers"
	"github.com/mohamedfawas/user-management-system/initializers"
	"github.com/mohamedfawas/user-management-system/utils"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	port := os.Getenv("PORT")

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:" + port}
	router.Use(cors.New(config))

	router.LoadHTMLGlob("templates/*")

	//routes: user login
	router.GET("/", handlers.DisplaySignInPage)
	router.POST("/", handlers.PostUserLogin)

	//routes: user sign up
	router.GET("/usersignup", handlers.DisplayUserSignUp)
	router.POST("/usersignup", handlers.PostUserSignUp)

	//routes : admin login
	router.GET("/admin", handlers.DisplayAdminLogin)
	router.POST("/admin", handlers.PostAdminLogin)

	//routes : admin sign up
	router.GET("/admin/signup", handlers.DisplayAdminSignUp)
	router.POST("/admin/signup", handlers.PostAdminSignUp)

	// from admin panel
	admin := router.Group("/admin")
	admin.Use(utils.AuthMiddleware(true))
	{
		admin.GET("/panel", controllers.DisplayAdminPanel)
		admin.GET("/panel/searchuser", controllers.SearchUser)
		admin.POST("/panel/createuser", controllers.PostCreateUser)
		admin.POST("/panel/deleteuser/:id", controllers.DeleteUser)
	}

	//set port and run
	router.Run(":" + port)
}
