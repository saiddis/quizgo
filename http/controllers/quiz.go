package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"gihub.com/saiddis/quizgo"
	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var quizzesUrl = []string{
	//"https://opentdb.com/api.php?amount=3&difficulty=easy",
	"https://opentdb.com/api.php?type=multiple&amount=3&difficulty=medium",
	//"https://opentdb.com/api.php?type=multiple&amount=2&difficulty=hard",
}

type TriviaCaller struct {
	client *http.Client
}

func NewTriviaCaller(client *http.Client) *TriviaCaller {
	return &TriviaCaller{
		client: client,
	}
}

type Trivia struct {
	ResponseCode int           `json:"response_code"`
	Results      []quizgo.Quiz `json:"results"`
}

type quizzes []quizgo.Quiz

func (t *TriviaCaller) getQuizzes(c *gin.Context, quizCategory string) (*quizzes, error) {
	if quizCategory != "" {
		for _, url := range quizzesUrl {
			url += "&" + "category=" + quizCategory
		}
	}

	//quizStore := quizgo.NewQuizStore()
	quizzes := quizzes{}
	quizGetter := t.callServer(&quizzes)
	difficulties := map[int]string{
		0: "medium",
		1: "hard",
	}

	var wg sync.WaitGroup
	wg.Add(len(quizzesUrl))

	for i, url := range quizzesUrl {
		//if i < 1 {
		//	time.Sleep(time.Second * 3)
		//}
		go func() {
			defer wg.Done()
			quizGetter(c, url, difficulties[i])
		}()
	}
	wg.Wait()

	if c.Errors != nil {
		return nil, errors.New("error getting quizzes")
	}
	return &quizzes, nil
}

func (t *TriviaCaller) callServer(quizzes *quizzes) func(c *gin.Context, url, difficulty string) {
	return func(c *gin.Context, url, difficulty string) {
		url += "&" + "difficulty=" + difficulty
		req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
		if err != nil {
			c.JSON(400, gin.H{"request error": err.Error()})
			return
		}
		resp, err := t.client.Do(req)
		if err != nil {
			c.JSON(400, gin.H{"response error": err.Error()})
			return
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(400, gin.H{"read error": err.Error()})
			return
		}
		defer resp.Body.Close()

		trivia := Trivia{}
		err = json.Unmarshal(data, &trivia)
		if err != nil {
			c.JSON(400, gin.H{"unmarshal error": err.Error()})
			return
		}

		if trivia.ResponseCode != 0 {
			c.JSON(400, gin.H{"error": "failed to fetch quiz data"})
			return
		}

		*quizzes = append(*quizzes, trivia.Results...)

		//quizStore.Mux.Lock()
		//quizStore.Mux.Unlock()
	}
}

func (s *Server) CreateQuiz(c *gin.Context) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	quizType := c.Request.FormValue("type")
	quizCategory := c.Request.FormValue("category")

	triviaCaller := NewTriviaCaller(client)
	quizzes, err := triviaCaller.getQuizzes(c, quizCategory)
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

	_, err = s.db.CreateQuiz(c, database.CreateQuizParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		QuizType:     quizType,
		QuizCategory: quizCategory,
		UserID:       uuid.MustParse(user.ID.String()),
	})
	if err != nil {
		log.Printf("couldn't create quiz: %v", err)
	}
	c.HTML(200, "quiz.html", gin.H{"quizzes": quizzes})
}
