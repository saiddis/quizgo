package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gihub.com/saiddis/quizgo/http/client"
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

	fmt.Println(dbName, dbPassword)
	db, err := postgres.NewDB(dbName,
		postgres.WithUser("saiddis"),
		postgres.WithPassword(dbPassword),
		postgres.WithSSL("disable"),
		postgres.WithTimeZone("Asia/Dushanbe"))
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	auth, err := authenticator.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	triviaCaller := client.NewTriviaCaller(&http.Client{
		Timeout: 30 * time.Second,
	})

	rtr := controllers.NewServer(db, triviaCaller, auth)
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
