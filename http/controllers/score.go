package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	if err = c.BindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error unmarshalling score params: %v", err)})
		return
	}
	log.Printf("%+v", params)

	session := sessions.Default(c)
	profile := session.Get("profile")
	var email string
	if profile, ok := profile.(map[string]interface{}); ok {
		email = profile["email"].(string)
	}

	user, err := s.db.GetUserByEmail(c, email)
	if err != nil {
		log.Printf("error retrieving user by email: %v", email)
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user by email: %v", err)})
		return
	}

	scoreID, err := s.db.CreateScore(c, database.CreateScoreParams{
		CompletionTime:    params.CompletionTime,
		HardQuizzesDone:   int32(params.HardQuizzesDone),
		MediumQuizzesDone: int32(params.MediumQuizzesDone),
		EasyQuizzesDone:   int32(params.EasyQuizzesDone),
		TotalScore:        int32(params.TotalScore),
		UserID:            user.ID,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error creating score: %v", err)})
		return
	}

	pgInt := pgtype.Int8{
		Int64: scoreID,
		Valid: true,
	}

	quizID, err := strconv.Atoi(params.QuizIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error converting quid_id from request: %v", err)})
		return
	}

	_, err = s.updateScoreID(c, int64(quizID), pgInt)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error updating score_id: %v", err)})
		log.Printf("failed to update score_id: %v", err)
		return
	}

	c.JSON(201, scoreID)
}

func (s *Server) GetScore(c *gin.Context) {
	scoreIDStr := c.Request.FormValue("id")
	scoreID, err := strconv.Atoi(scoreIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error converting score_id: %v", err)})
		log.Printf("couldn't convert score_id: %v", err)
		return
	}
	score, err := s.db.GetScoreByID(c, int64(scoreID))
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error retrieving score by id: %v", err)})
		log.Printf("failed to get score by id: %v", err)
		return
	}

	c.JSON(200, score)

}

func (s *Server) updateScoreID(c *gin.Context, quizID int64, scoreID pgtype.Int8) (database.Quiz, error) {
	quiz, err := s.db.UpdateScoreID(c, database.UpdateScoreIDParams{
		ScoreID: scoreID,
		ID:      quizID,
	})
	if err != nil {
		return database.Quiz{}, fmt.Errorf("error updating score_id in quiz: %v", err)
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

func (s *Server) GetHighestUserScore(c *gin.Context) {
	email := c.Request.FormValue("email")
	score, err := s.db.GetUserHighestScoreByEmail(c, email)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error retrieving highest user score score by id: %v", err)})
		log.Printf("error retrieving highest user score by email: %v", err)
		return
	}

	c.JSON(200, score)
}

func extractCompletionTime(m map[string]interface{}) (time.Duration, error) {
	var completionTime time.Duration
	var err error
	for k, v := range m {
		if k == "completion_time" {
			if v, ok := v.(string); ok {
				completionTime, err = time.ParseDuration(v)
				if err != nil {
					return 0, fmt.Errorf("error converting completion_time to time.Duration: %v", err)
				}
			} else {
				return 0, fmt.Errorf("completion_time type must be string not %v", v)
			}
			break
		}
	}

	return completionTime, nil
}
