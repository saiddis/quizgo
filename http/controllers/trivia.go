package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetTrivias(c *gin.Context) {
	quizIDStr := c.Request.FormValue("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("coldn't convert quiz_id: %v", err)})
		return
	}

	trivias, err := s.db.GetTriviasByQuizID(c, int64(quizID))
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("coldn't get trivias by quiz_id: %v", err)})
		return
	}

	c.JSON(200, trivias)
}
