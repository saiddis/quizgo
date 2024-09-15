package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	_, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database:", err)
	}

}
