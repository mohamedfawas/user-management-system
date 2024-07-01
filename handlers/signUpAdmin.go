package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-management-system/initializers"
	"github.com/mohamedfawas/user-management-system/models"
	"golang.org/x/crypto/bcrypt"
)

func DisplayAdminSignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "adminSignUp.html", gin.H{
		"title": "admin sign up",
	})
}

func PostAdminSignUp(c *gin.Context) {
	var admin models.Admin
	admin.UserName = c.PostForm("username")
	admin.Email = c.PostForm("email")
	admin.Password = c.PostForm("password")

	c.Bind(&admin)
	// //If email already exists
	var existingAdmin models.Admin
	if err := initializers.DB.Where("email = ?", admin.Email).First(&existingAdmin).Error; err == nil {
		c.HTML(http.StatusConflict, "adminSignUp.html", gin.H{
			"error":   "Given email already exists",
			"message": "Create the account with a different email",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	admin.Password = string(hashedPassword)
	if err := initializers.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	c.HTML(http.StatusOK, "adminLogin.html", nil)
}
