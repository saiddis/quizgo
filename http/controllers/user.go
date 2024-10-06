package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for our logged-in user page.
func UserHandler(c *gin.Context) {
	session := sessions.Default(c)
	profile := session.Get("profile")

	c.HTML(http.StatusOK, "user.html", gin.H{
		"profile": profile,
	})
}
