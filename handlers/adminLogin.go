package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-management-system/initializers"
	"github.com/mohamedfawas/user-management-system/models"
	"github.com/mohamedfawas/user-management-system/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func DisplayAdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminLogin.html", gin.H{
		"title": "Welcome to admin login page",
	})
}

func PostAdminLogin(c *gin.Context) {
	var admin models.Admin
	//admin.Name = c.PostForm("name")
	admin.Email = c.PostForm("email")
	admin.Password = c.PostForm("password")

	if err := c.ShouldBind(&admin); err != nil {
		c.HTML(http.StatusBadRequest, "adminLogin.html", gin.H{"error": "Invalid input"})
		return
	}

	var dbAdmin models.Admin
	result := initializers.DB.Where("email = ?", admin.Email).First(&dbAdmin)
	if result.Error != nil {
		log.Printf("Database query error: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			c.HTML(http.StatusUnauthorized, "adminLogin.html", gin.H{"error": "Admin not found"})
		} else {
			c.HTML(http.StatusInternalServerError, "adminLogin.html", gin.H{"error": "Database error"})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "adminLogin.html", gin.H{"error": "Invalid login credentials"})
		return
	}

	token, err := utils.GenerateToken(dbAdmin.ID, true)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "adminLogin.html", gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("jwt", token, int(time.Hour*24/time.Second), "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/panel")
}
