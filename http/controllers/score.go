package controllers

import (
	"fmt"
	"log"
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) CreateScore(c *gin.Context) {
	var params struct {
		HardQuizzesDone   int `json:"hard_quizzes_done"`
		MediumQuizzesDone int `json:"medium_quizzes_done"`
		EasyQuizzesDone   int `json:"easy_quizzes_done"`
		TotalScore        int `json:"total_score"`
	}
	rawParams := make(map[string]interface{})
	var err error

	if err = c.BindJSON(&rawParams); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to parse json: %v", err)})
		return
	}

	if err = c.BindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to parse json: %v", err)})
		return
	}

	completionTime, err := extractCompletionTime(rawParams)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	profile := session.Get("profile")
	var email string
	if profile, ok := profile.(map[string]interface{}); ok {
		email = profile["email"].(string)
	}

	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		log.Printf("couldn't get user by email: %v", email)
	}

	score, err := s.db.CreateScore(c, database.CreateScoreParams{
		ID:                uuid.New(),
		CompletionTime:    completionTime.Milliseconds(),
		HardQuizzesDone:   int32(params.HardQuizzesDone),
		MediumQuizzesDone: int32(params.MediumQuizzesDone),
		EasyQuizzesDone:   int32(params.EasyQuizzesDone),
		TotalScore:        int32(params.TotalScore),
		UserID:            user.ID,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to create score: %v", err)})
		return
	}
	c.JSON(201, score)
}

func extractCompletionTime(m map[string]interface{}) (time.Duration, error) {
	var completionTime time.Duration
	var err error
	for k, v := range m {
		if k == "completion_time" {
			if v, ok := v.(string); ok {
				completionTime, err = time.ParseDuration(v)
				if err != nil {
					return 0, fmt.Errorf("Error converting completion_time to time.Duration: %v", err)
				}
			} else {
				return 0, fmt.Errorf("completion_time type must be string not %v", v)
			}
			break
		}
	}

	return completionTime, nil
}
