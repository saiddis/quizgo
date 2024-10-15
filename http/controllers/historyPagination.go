package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetLastQuizIDByEmail(c *gin.Context) {
	email := c.Request.FormValue("email")
	if strings.Contains(email, "%40") {
		email = strings.Replace(email, "%40", "@", 1)
	}

	log.Println(email)
	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id by email: %v", err)})
		log.Printf("error retrieving user_id by email: %v", err)
		return
	}

	quiz, err := s.db.GetLastQuizByUserID(c, user.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving quiz_id by user_id: %v", err)})
		log.Printf("error retrieving quiz_id by user_id: %v", err)
		return
	}

	c.HTML(200, "history.html", gin.H{"quiz_id": quiz.ID})
}

func (s *Server) QuizzesPagination(c *gin.Context) {
	//var params struct {
	//	LastQuizID int `json:"quiz_id"`
	//}
	//err := c.BindJSON(&params)
	//if err != nil {
	//	c.JSON(400, gin.H{"error": fmt.Errorf("error unmarshalling quiz_id from context: %v", err)})
	//	return
	//}
	lastQuizIDStr := c.Request.FormValue("id")
	lastQuizID, err := strconv.Atoi(lastQuizIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting quiz_id: %v", err)})
		return
	}
	userID, err := s.GetUserIDByContext(c)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id from context: %v", err)})
		return
	}

	quizzes, err := s.db.QuizzesPagination(c, database.QuizzesPaginationParams{
		UserID: userID,
		ID:     int64(lastQuizID),
	})
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving history data: %v", err)})
		return
	}

	c.JSON(200, quizzes)
}
