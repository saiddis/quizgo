package controllers

import (
	"fmt"
	"log"
	"time"

	"gihub.com/saiddis/quizgo"
	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var quizzesUrl = []string{
	//"https://opentdb.com/api.php?amount=3&difficulty=easy",
	"https://opentdb.com/api.php?amount=5",
	//"https://opentdb.com/api.php?type=multiple&amount=2&difficulty=hard",
}

func (s *Server) CreateQuizForGuest(c *gin.Context) {
	quizType := c.Request.FormValue("type")
	quizCategory := c.Request.FormValue("category")

	urls := getAddedURLParams(quizType, quizCategory)

	trivias, err := s.Client.Fetch(c, urls)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "quiz.html", gin.H{"quizzes": trivias, "quiz_id": ""})

}

func (s *Server) CreateQuizForUser(c *gin.Context) {
	quizType := c.Request.FormValue("type")
	quizCategory := c.Request.FormValue("category")

	urls := getAddedURLParams(quizType, quizCategory)

	trivias, err := s.Client.Fetch(c, urls)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, err := s.GetUserIDByContext(c)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error retrieving user_id from context: %v", err)})
		return
	}
	quizID, err := s.db.CreateQuiz(c, database.CreateQuizParams{
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Type:      quizType,
		Category:  quizCategory,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error creating quiz: %v", err)})
		return
	}

	*trivias, err = s.insertHistoryData(c, *trivias, quizID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("error inserting history data: %v", err)})
		log.Printf("error inserting history data: %v", err)
		return
	}

	c.HTML(201, "quiz.html", gin.H{"quizzes": *trivias, "quiz_id": quizID})
}

func (s *Server) insertHistoryData(c *gin.Context, trivias []quizgo.Trivia, quizID int64) ([]quizgo.Trivia, error) {
	var err error
	var triviaID uuid.UUID
	var dbTrivia database.Trivia
	var quizzesTriviasParams []database.CreateQuizzesTriviasParams
	var triviasParams []database.CreateTriviasParams
	var optionsParams []database.CreateOptionsParams
	var optionID uuid.UUID

	for i := 0; i < len(trivias); i++ {
		dbTrivia, err = s.db.GetTriviaByQuestion(c, trivias[i].Question)
		if err == nil {
			triviaID = dbTrivia.ID
			trivias[i].ID = triviaID
			options, err := s.db.GetOptionsIDByTriviaID(c, triviaID)
			if err != nil {
				return nil, fmt.Errorf("error retrieving options by trivia_id: %v", err)
			}
			for _, row := range options {
				if !row.Correct {
					trivias[i].OptionsID = append(trivias[i].OptionsID, row.ID)
				} else {
					trivias[i].CorrectOptionID = row.ID
				}
			}
		} else {
			triviaID = uuid.New()
			trivias[i].ID = triviaID
			log.Printf("Inserting trivia: %s", trivias[i].Question)
			triviasParams = append(triviasParams, database.CreateTriviasParams{
				ID:         triviaID,
				Type:       trivias[i].Type,
				Category:   trivias[i].Category,
				Difficulty: trivias[i].Difficulty,
				Question:   trivias[i].Question,
			})
			for _, option := range trivias[i].Options {
				optionID = uuid.New()
				optionsParams = append(optionsParams, database.CreateOptionsParams{
					ID:       optionID,
					Option:   option,
					Correct:  false,
					TriviaID: triviaID,
				})
				trivias[i].OptionsID = append(trivias[i].OptionsID, optionID)
			}
			optionID = uuid.New()
			optionsParams = append(optionsParams, database.CreateOptionsParams{
				ID:       optionID,
				Option:   trivias[i].CorrectOption,
				Correct:  true,
				TriviaID: triviaID,
			})
			trivias[i].CorrectOptionID = optionID

			quizzesTriviasParams = append(quizzesTriviasParams, database.CreateQuizzesTriviasParams{
				QuizID:   quizID,
				TriviaID: triviaID,
			})
		}
	}

	var n int64
	if triviasParams != nil {
		log.Printf("Attempting to insert %d trivia records", len(triviasParams))
		n, err = s.db.CreateTrivias(c, triviasParams)
		if err != nil {
			return nil, fmt.Errorf("error bulk inserting trivias: %v", err)
		}
		log.Printf("%d trivias inserted into trivias table", n)

		log.Printf("Attempting to insert %d options", len(optionsParams))
		n, err = s.db.CreateOptions(c, optionsParams)
		if err != nil {
			return nil, fmt.Errorf("error bulk inserting options: %v", err)
		}
		log.Printf("%d options inserted into options table", n)

	}

	n, err = s.db.CreateQuizzesTrivias(c, quizzesTriviasParams)
	if err != nil {
		return nil, fmt.Errorf("error creating quiz_trivia: %v", err)
	}
	log.Printf("%d quizzes_trivias inserted into quizzes_trivias table", n)

	return trivias, nil
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
