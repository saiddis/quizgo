package controllers

import (
	"encoding/gob"
	"html/template"
	"io"
	"log"

	"gihub.com/saiddis/quizgo"
	"gihub.com/saiddis/quizgo/internal/install/database"
	"gihub.com/saiddis/quizgo/middleware/authenticator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Client interface {
	Fetch(c *gin.Context, urls []string) (*[]quizgo.Trivia, error)
}

type Server struct {
	Router *gin.Engine
	db     *database.Queries
	auth   *authenticator.Authenticator
	Client Client
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c *gin.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("../../web/template/*.html")),
	}
}
func NewServer(db *database.Queries, client Client, auth *authenticator.Authenticator) *Server {
	server := &Server{
		Router: gin.Default(),
		db:     db,
		auth:   auth,
		Client: client,
	}

	html := newTemplate()
	server.Router.SetHTMLTemplate(html.templates)
	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	env, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatalf("Could't read .env file: %v", err)
	}
	secret := env["AUTH0_CLIENT_SECRET"]
	store := cookie.NewStore([]byte(secret))
	server.Router.Use(sessions.Sessions("auth-session", store))

	server.Router.Static("/css", "../../web/css")
	server.Router.Static("/static", "../../web/static")

	server.Router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home.html", nil)
	})
	server.Router.GET("/login", LoginHandler(server.auth))
	server.Router.GET("/callback", CallbackHandler(server.auth))
	server.Router.GET("/logout", LogoutHandler)
	server.Router.GET("/quiz", server.CreateQuizForGuest)

	server.Router.GET("/leaderboard", server.GetTheHighestScore)
	server.Router.GET("/leaderboard/load", server.UsersBestScorePagination)

	user := server.Router.Group("/user")
	user.Use(auth.IsAuthenticated)
	user.GET("/", server.UserHandler)
	user.GET("/quiz", server.CreateQuizForUser)
	user.POST("/quiz/score", server.CreateScore)
	user.POST("/quiz/answer", server.CreateAnswer)

	user.GET("/history", server.GetLastQuizIDByEmail)
	user.GET("/history/load", server.QuizzesPagination)
	user.GET("/history/score", server.GetScore)
	user.GET("/history/trivia", server.GetTrivias)
	user.GET("/history/answer", server.GetAnswersByQuizID)
	user.GET("/history/option", server.GetOptionByAnswerID)
	return server
}
