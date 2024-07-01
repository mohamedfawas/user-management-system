package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-management-system/initializers"
	"github.com/mohamedfawas/user-management-system/models"
	"golang.org/x/crypto/bcrypt"
)

func DisplayCreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "createUserByAdmin.html", gin.H{
		"title": "welcome to create user page",
	})
}

func PostCreateUser(c *gin.Context) {
	var user models.User
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	password := c.PostForm("password")

	// Check if email already exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		DisplayAdminPanelWithMessage(c, "Email already exists", "error")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		DisplayAdminPanelWithMessage(c, "Failed to hash password", "error")
		return
	}
	user.Password = string(hashedPassword)

	// Create the user
	if err := initializers.DB.Create(&user).Error; err != nil {
		DisplayAdminPanelWithMessage(c, "Failed to create user", "error")
		return
	}

	DisplayAdminPanelWithMessage(c, "New user created successfully", "message")
}

func DisplayAdminPanelWithMessage(c *gin.Context, message string, messageType string) {
	var users []models.User
	if err := initializers.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "adminPanel.html", gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	c.HTML(http.StatusOK, "adminPanel.html", gin.H{
		"users":     users,
		messageType: message,
	})
}

func DisplayAdminPanel(c *gin.Context) {
	DisplayAdminPanelWithMessage(c, "", "")
}

func SearchUser(c *gin.Context) {
	searchTerm := c.Query("searchTerm")
	var users []models.User

	if searchTerm == "" {
		c.Redirect(http.StatusSeeOther, "/admin/panel")
		return
	}

	// First, search by name and email
	result := initializers.DB.Where("name LIKE ? OR email LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").Find(&users)

	// If no results and the searchTerm is numeric, try searching by ID
	if len(users) == 0 {
		if id, err := strconv.ParseUint(strings.TrimSpace(searchTerm), 10, 64); err == nil {
			initializers.DB.Where("id = ?", id).Find(&users)
		}
	}

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "adminPanel.html", gin.H{
			"error": "Failed to search users: " + result.Error.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "adminPanel.html", gin.H{
		"users":      users,
		"searchTerm": searchTerm,
	})
}

func DeleteUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Param("id")

	// Delete the user from the database
	result := initializers.DB.Delete(&models.User{}, userID)
	if result.Error != nil {
		DisplayAdminPanelWithMessage(c, "Failed to delete user", "error")
		return
	}

	if result.RowsAffected == 0 {
		DisplayAdminPanelWithMessage(c, "User not found", "error")
		return
	}

	// Display success message
	message := fmt.Sprintf("User with ID %s has been deleted", userID)
	DisplayAdminPanelWithMessage(c, message, "message")
}
