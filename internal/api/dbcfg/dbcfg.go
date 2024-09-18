package dbcfg

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gihub.com/saiddis/quizgo/internal/install/database"
)

type options struct {
	url      string
	protocol string
	user     string
	host     string
	password string
	port     *int
}

type Option func(opts *options) error

func WithURL(url string) Option {
	return func(opts *options) error {
		opts.url = url
		return nil
	}
}

func connectToDB(db, url string) (*database.Queries, error) {
	conn, err := sql.Open(db, url)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	return database.New(conn), nil
}

func NewDB(db, dbName string, opts ...Option) (*database.Queries, error) {
	var options options
	var err error
	for _, opt := range opts {
		err = opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var url string
	if options.url != "" {
		url = options.url
		return connectToDB(db, url)
	}
	protocol := "postgres"
	user := "postgres"
	port := 5432
	password := "root"
	host := "localhost"
	sslMode := "?sslmode=disable"

	if options.protocol != "" {
		protocol = options.protocol
	}
	if options.user != "" {
		user = options.user
	}
	if options.port != nil {
		port = *options.port
		if port <= 0 {
			return nil, errors.New("port parameter has to be more than 0")
		}
	}
	if options.password != "" {
		password = options.password
	}

	url = fmt.Sprintf("%s://%s:%s@%s:%d/%s%s", protocol, user, password, host, port, dbName, sslMode)
	return connectToDB(dbName, url)
}
