package controllers

import (
	"log"
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var quizzesUrl = []string{
	//"https://opentdb.com/api.php?amount=3&difficulty=easy",
	"https://opentdb.com/api.php?amount=3",
	//"https://opentdb.com/api.php?type=multiple&amount=2&difficulty=hard",
}

func (s *Server) CreateQuiz(c *gin.Context) {
	quizType := c.Request.FormValue("type")
	quizCategory := c.Request.FormValue("category")

	urls := getAddedURLParams(quizType, quizCategory)

	quizzes, err := s.Client.Fetch(c, urls)
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

	quiz, err := s.db.CreateQuiz(c, database.CreateQuizParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		QuizType:     quizType,
		QuizCategory: quizCategory,
		UserID:       uuid.MustParse(user.ID.String()),
	})
	if err != nil {
		log.Printf("couldn't create quiz: %v", err)
	}
	c.HTML(200, "quiz.html", gin.H{"quizzes": quizzes, "quiz_id": quiz.ID})
}

func getAddedURLParams(quizType, quizCategory string) []string {
	urls := quizzesUrl
	if quizCategory != "" && quizCategory != "any" {
		for i := range urls {
			urls[i] += "&" + "category=" + quizCategory
		}
	}
	if quizType != "" && quizType != "any" {
		for i := range urls {
			urls[i] += "&" + "type=" + quizType
		}
	}
	return urls
}
