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
		CompletionTime    int64  `json:"completion_time"`
		HardQuizzesDone   int    `json:"hard_quizzes_done"`
		MediumQuizzesDone int    `json:"medium_quizzes_done"`
		EasyQuizzesDone   int    `json:"easy_quizzes_done"`
		TotalScore        int    `json:"total_score"`
		QuizIDStr         string `json:"quiz_id"`
	}
	var err error

	//rawParams := make(map[string]interface{})
	//if err = c.BindJSON(&rawParams); err != nil {
	//	c.JSON(400, gin.H{"error": fmt.Errorf("failed to parse json: %v", err)})
	//	return
	//}

	if err = c.BindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to parse json: %v", err)})
		return
	}

	log.Printf("%+v", params)

	//completionTime, err := extractCompletionTime(rawParams)
	//if err != nil {
	//	c.JSON(400, gin.H{"error": fmt.Errorf("couldn't extract completion time: %v", err)})
	//	return
	//}

	session := sessions.Default(c)
	profile := session.Get("profile")
	var email string
	if profile, ok := profile.(map[string]interface{}); ok {
		email = profile["email"].(string)
	}

	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		log.Printf("couldn't get user by email: %v", email)
		c.JSON(400, gin.H{"error": fmt.Errorf("couldn't get user by email: %v", err)})
		return
	}

	//t, err := time.ParseDuration(params.CompletionTime)
	//if err != nil {
	//	log.Printf("couldn't parse duration: %v", email)
	//	c.JSON(400, gin.H{"error": fmt.Errorf("couldn't parse duration: %v", err)})
	//	return
	//}

	score, err := s.db.CreateScore(c, database.CreateScoreParams{
		ID:                uuid.New(),
		CompletionTime:    params.CompletionTime,
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

	_, err = s.updateScoreID(c, params.QuizIDStr, score.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to update score_id: %v", err)})
		log.Printf("failed to update score_id: %v", err)
		return
	}

	c.JSON(201, score)
}

func (s *Server) updateScoreID(c *gin.Context, quizIDStr string, scoreID uuid.UUID) (database.Quiz, error) {
	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		return database.Quiz{}, fmt.Errorf("failed to parse %s to uuid: %v", quizIDStr, err)
	}
	quiz, err := s.db.UpdateScoreID(c, database.UpdateScoreIDParams{
		ScoreID: s.newNullUUID(scoreID),
		ID:      quizID,
	})
	if err != nil {
		return database.Quiz{}, fmt.Errorf("failed to update score_id in quiz: %v", err)
	}
	return quiz, nil
}

func (s *Server) newNullUUID(id uuid.UUID) uuid.NullUUID {
	if len(id) == 0 {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  id,
		Valid: true,
	}
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
