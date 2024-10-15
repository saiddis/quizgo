package controllers

import (
	"fmt"
	"strings"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetLastQuizIDByEmail(c *gin.Context) {
	email := c.Request.FormValue("email")
	if strings.Contains(email, "%40") {
		email = strings.Replace(email, "%40", "@", 1)
	}
	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id by email: %v", err)})
		return
	}

	quiz, err := s.db.GetLastQuizByUserID(c, user.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id from conetext: %v", err)})
		return
	}

	c.HTML(200, "history.html", gin.H{"quiz_id": quiz.ID})
}

func (s *Server) QuizzesPagination(c *gin.Context) {
	var params struct {
		LastQuizID int `json:"quiz_id"`
	}
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error unmarshalling quiz_id from context: %v", err)})
		return
	}
	userID, err := s.GetUserIDByContext(c)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id from context: %v", err)})
		return
	}

	quizzes, err := s.db.PaginateQuizzes(c, database.PaginateQuizzesParams{
		UserID: userID,
		ID:     int64(params.LastQuizID),
	})
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error loading history data: %v", err)})
		return
	}

	c.JSON(200, quizzes)
}
