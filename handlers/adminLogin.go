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

func DisplayAdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminLogin.html", gin.H{
		"title": "Welcome to admin login page",
	})
}

func PostAdminLogin(c *gin.Context) {
	var admin models.Admin
	admin.Email = c.PostForm("email")
	admin.Password = c.PostForm("password")

	if err := c.Bind(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var dbAdmin models.Admin
	if err := initializers.DB.Where("email = ?", admin.Email).First(&dbAdmin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	token, err := utils.GenerateToken(dbAdmin.ID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("jwt", token, int(time.Hour*24/time.Second), "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/panel")
}
