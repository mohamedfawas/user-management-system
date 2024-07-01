package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-management-system/initializers"
	"github.com/mohamedfawas/user-management-system/models"
	"github.com/mohamedfawas/user-management-system/utils"
	"golang.org/x/crypto/bcrypt"
)

func DisplayUserSignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "userSignUp.html", gin.H{
		"title": "User Sign Up ",
	})
}

func PostUserSignUp(c *gin.Context) {
	var user models.User
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	c.Bind(&user) //bind data from an HTTP request body to a Go struct
	//Form parsing vs JSON binding: Your HTML form is sending data as form-urlencoded , so here we use "bind" instead of "bindjson"

	//If email already exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusConflict, "userSignUp.html", gin.H{
			"error":   "Given email already exists",
			"message": "Create the account with a different email",
		})
		return
	}

	//hash the given password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "userSignUp.html", gin.H{
			"error": "Failed to hash password",
		})
		return
		//add additional options after testing
	}
	user.Password = string(hashedPassword)
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "userSignUp.html", gin.H{
			"error": "Failed to create user",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, false)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "userSignUp.html", gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.SetCookie("jwt", token, int(time.Hour*24/time.Second), "/", "localhost", false, true)
	c.HTML(http.StatusOK, "index.html", nil)
}
