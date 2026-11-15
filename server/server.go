package server

import (
	"MokutonCoffee-API/config"
	"MokutonCoffee-API/controllers"
	"MokutonCoffee-API/middleware"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// New membangun gin engine lengkap dengan middleware dan routes.
// Dipakai bersama oleh mode lokal (main.go) dan serverless Vercel (api/index.go).
func New() *gin.Engine {
	// Di Vercel env datang dari dashboard; .env hanya relevan saat lokal.
	_ = godotenv.Load()

	if os.Getenv("VERCEL") != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.ConnectDB()

	s := gin.Default()
	s.SetTrustedProxies(nil)

	// jangan "*": dengan AllowCredentials: true browser menolak wildcard
	origins := []string{}
	for _, o := range strings.Split(os.Getenv("CLIENT_ORIGIN"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}

	s.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	api := s.Group("/api")
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

	return s
}
