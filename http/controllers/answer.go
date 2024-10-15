package controllers

import (
	"fmt"
	"strconv"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) CreateAnswer(c *gin.Context) {
	var params struct {
		QuizIDStr   string `json:"quiz_id"`
		TriviaIDStr string `json:"trivia_id"`
		OptionIDStr string `json:"option_id"`
	}

	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error unmarshalling request: %v", err)})
		return
	}

	optionID, err := uuid.Parse(params.OptionIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting option_id: %v", err)})
		return
	}

	quizID, err := strconv.Atoi(params.QuizIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting quiz_id: %v", err)})
		return
	}
	triviaID, err := uuid.Parse(params.TriviaIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error convert trivia_id: %v", err)})
		return
	}

	answerID, err := s.db.CreateAnswer(c, database.CreateAnswerParams{
		QuizID:   int64(quizID),
		TriviaID: triviaID,
		OptionID: optionID,
	})

	c.JSON(201, answerID)
}

func (s *Server) GetAnswersByQuizID(c *gin.Context) {
	quizIDStr := c.Request.FormValue("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting quiz_id: %v", err)})
		return
	}

	answers, err := s.db.GetAnswersByQuizID(c, int64(quizID))
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving answers by quiz_id: %v", err)})
		return
	}

	c.JSON(200, answers)
}

func (s *Server) GetOptionByAnswerID(c *gin.Context) {
	answerIDStr := c.Request.FormValue("id")
	answerID, err := strconv.Atoi(answerIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting answer_id: %v", err)})
		return
	}

	option, err := s.db.GetOptionByAnswerID(c, int64(answerID))
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving option by answer_id: %v", err)})
		return
	}

	c.JSON(200, option)
}
