package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserIDByContext(c *gin.Context) (uuid.UUID, error) {
	session := sessions.Default(c)
	var email string

	if profile := session.Get("profile"); profile != nil {
		if profile, ok := profile.(map[string]interface{}); ok {
			email = profile["email"].(string)
		} else {
			return uuid.UUID{}, fmt.Errorf("error retrieving email from profile: %v", profile)
		}
	} else {
		return uuid.UUID{}, fmt.Errorf("no profile data found in context: %v", c)
	}

	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error retrieving user by email: %v", err)
	}

	return user.ID, nil
}

func (s *Server) GetUserIDByQuizID(c *gin.Context) {
	var params struct {
		QuizID string `json:"quiz_id"`
	}

	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error unmarshalling quiz_id from context: %v", err)})
		return
	}

	quizID, err := strconv.Atoi(params.QuizID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting quiz_id from context: %v", err)})
		return
	}

	userID, err := s.db.GetUserIDByQuizID(c, int64(quizID))

	c.JSON(200, gin.H{"user_id": userID})
}

// Handler for our logged-in user page.
func UserHandler(c *gin.Context) {
	session := sessions.Default(c)
	profile := session.Get("profile")
	var pictureSrc string
	var userEmail string
	if profile, ok := profile.(map[string]interface{}); ok {
		if picture, ok := profile["picture"].(string); ok {
			pictureSrc = picture
		}
		if email, ok := profile["email"].(string); ok {
			userEmail = email
		}
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"profile": profile,
		"picture": pictureSrc,
		"email":   userEmail,
	})
}
