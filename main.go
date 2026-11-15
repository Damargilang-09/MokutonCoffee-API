package main

import (
	"MokutonCoffee-API/config"
	"MokutonCoffee-API/controllers"
	"MokutonCoffee-API/middleware"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file")
	}

	config.ConnectDB()

	server := gin.Default()
	origins := strings.Split(os.Getenv("CLIENT_ORIGIN"), ",")
	server.Use(cors.New(cors.Config{
		AllowOrigins:     origins, // jangan "*": dengan AllowCredentials: true browser menolak "*"
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Route
	api := server.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("", middleware.Authentication(), middleware.Authorization("admin"), controllers.CreateProduct)
			products.GET("", controllers.GetAllProducts)
			products.GET("/:id", controllers.GetProductById)
			products.PATCH("/:id", middleware.Authentication(), middleware.Authorization("admin"), controllers.UpdateProduct)
			products.DELETE("/:id", middleware.Authentication(), middleware.Authorization("admin"), controllers.DeleteProduct)
		}

		user := api.Group("/auth")
		{
			user.POST("/register", controllers.RegisterUser)
			user.POST("/login", controllers.LoginUser)
			user.GET("/me", middleware.Authentication(), controllers.GetCurrentUser)
			user.POST("/logout", controllers.LogoutUser)
		}

	}

	server.Run(":8000")

}
