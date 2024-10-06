package main

import (
	"fmt"
	"log"

	"gihub.com/saiddis/quizgo/http/controllers"
	"gihub.com/saiddis/quizgo/middleware/authenticator"
	"gihub.com/saiddis/quizgo/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	env, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatalf("Could't read .env file: %v", err)
	}
	dbName := env["DB_NAME"]
	dbPassword := env["DB_PASSWORD"]

	// secret := env["AUTH0_CLIENT_SECRET"]
	fmt.Println(dbName, dbPassword)
	db, err := postgres.NewDB(dbName,
		postgres.WithUser("saiddis"),
		postgres.WithPassword(dbPassword),
		postgres.WithSSL("disable"),
		postgres.WithTimeZone("Asia/Dushanbe"))
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// tokenService := token.NewService(secret)

	// userService := user.NewService(db, tokenService)
	// scoreService := score.NewService(db)
	// sessionService := session.NewService(db)

	// router := chi.NewRouter()
	// router.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"*"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// }))

	// router.Post("/users", userService.CreateUser)
	// router.Get("/users", userService.GetUsers)
	// router.Post("/scores", scoreService.CreateScore)
	// router.Post("/sessions", sessionService.CreateSession)

	auth, err := authenticator.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := controllers.NewServer(db, auth)
	rtr.Router.Run()

	// srv, err := server.NewServer("localhost",
	// 	server.WithPort(8080),
	// 	server.WithHandler(rtr))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Server started on port %s", "8080")
	// err = srv.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
