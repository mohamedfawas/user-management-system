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

func DisplaySignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func PostUserLogin(c *gin.Context) {
	var user models.User
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input format, not able to bind the data",
		})
		return
	}

	var dbUser models.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid login credentials",
		})
		return
	}

	token, err := utils.GenerateToken(dbUser.ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.SetCookie("jwt", token, int(time.Hour*24/time.Second), "/", "localhost", false, true)

	c.HTML(http.StatusOK, "homePageUser.html", gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}
