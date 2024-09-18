package main

import (
	"log"
	"os"

	"gihub.com/saiddis/quizgo/internal/api/dbcfg"
	"gihub.com/saiddis/quizgo/internal/api/servcfg"
	"gihub.com/saiddis/quizgo/internal/handler"
	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	//	conn, err := sql.Open("postgres", dbURL)
	//	if err != nil {
	//		log.Fatal("Can't connect to database:", err)
	//	}
	//	db := database.New(conn)
	db, err := dbcfg.NewDB("postgres", "quizgo", dbcfg.WithURL(dbURL))
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	_ = apiConfig{
		DB: db,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/health", handler.Readiness)

	srv, err := servcfg.NewServer("localhost",
		servcfg.WithPort(8080),
		servcfg.WithHandler(router))
	if err != nil {

		log.Fatal(err)
	}

	log.Printf("Server started on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
