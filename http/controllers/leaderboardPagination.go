package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetTheHighestScore(c *gin.Context) {
	score, err := s.db.GetTheHighestTotalScore(c)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id from conetext: %v", err)})
		log.Printf("error retrieving user_id from context: %v", err)
		return
	}

	c.HTML(200, "leaderboard.html", gin.H{"score": score})
}

func (s *Server) UsersBestScorePagination(c *gin.Context) {
	highestScoreStr := c.Request.FormValue("score")
	highestScore, err := strconv.Atoi(highestScoreStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting score: %v", err)})
		log.Printf("error converting score: %v", err)
		return
	}

	rows, err := s.db.UsersBestScorePagination(c, int32(highestScore))
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving users best scores: %v", err)})
		log.Printf("error retrieving users best scores: %v", err)
		return
	}

	c.JSON(200, rows)
}
