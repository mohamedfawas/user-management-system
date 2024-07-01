package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Clear the JWT cookie
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)

	// Redirect to the login page
	c.Redirect(http.StatusSeeOther, "/")
}
