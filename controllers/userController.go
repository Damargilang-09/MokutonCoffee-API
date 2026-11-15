package controllers

import (
	"MokutonCoffee-API/config"
	"MokutonCoffee-API/models"
	"MokutonCoffee-API/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// setAuthCookie: di produksi frontend & backend berbeda subdomain .vercel.app
// (= cross-site karena vercel.app public suffix), sehingga wajib
// SameSite=None + Secure. Lokal tetap Lax + Secure=false.
func setAuthCookie(c *gin.Context, value string, maxAge int) {
	secure := os.Getenv("VERCEL") != ""
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	c.SetSameSite(sameSite)
	c.SetCookie("token", value, maxAge, "/", "", secure, true)
}

func RegisterUser(c *gin.Context) {
	var input AuthInputRegister

	//validation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Hash password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Sava to database
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error

	if userCreated != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to register account",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"Name":  user.Name,
			"Email": user.Email,
		},
	})
}

func LoginUser(c *gin.Context) {
	var input AuthInputLogin

	//validation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Email check
	var user models.User
	userData := config.DB.Where("email=?", input.Email).First(&user).Error
	if userData != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect email or password",
		})
		return
	}

	//password check
	errMatchPassword := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if errMatchPassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect email or password",
		})
		return
	}

	//Generate Token
	tokenString, err := utils.GenerateToken(user.ID.String(), string(user.Role))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed creating token",
		})
		return
	}

	setAuthCookie(c, tokenString, 3600*24)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"user": gin.H{
			"Name":  user.Name,
			"Email": user.Email,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	//get userId
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	var user models.User

	userData := config.DB.Select("id", "role", "name", "email").First(&user, userId).Error
	if userData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func LogoutUser(c *gin.Context) {
	setAuthCookie(c, "", -1)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successfully",
	})
}
