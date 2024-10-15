package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type parameters struct {
	sslMode  string
	timeZone string
}

type options struct {
	url      string
	user     string
	host     string
	password string
	port     int
	params   map[string]string
}

type Option func(opts *options)

func WithURL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

func WithHost(host string) Option {
	return func(opts *options) {
		opts.host = host
	}
}

func WithUser(usr string) Option {
	return func(opts *options) {
		opts.user = usr
	}
}

func WithPassword(pswd string) Option {
	return func(opts *options) {
		opts.password = pswd
	}
}

func WithPort(port int) Option {
	return func(opts *options) {
		opts.port = port
	}
}

func WithSSL(ssl string) Option {
	return func(opts *options) {
		opts.params["sslmode"] = ssl
	}
}

func WithTimeZone(tz string) Option {
	return func(opts *options) {
		opts.params["TimeZone"] = tz
	}
}

func NewDB(dbName string, opts ...Option) (*database.Queries, error) {
	options := options{
		user:     "postgres",
		host:     "localhost",
		password: "root",
		port:     5432,
		params:   map[string]string{},
	}
	for _, opt := range opts {
		opt(&options)
	}

	var url string
	if options.url != "" {
		url = options.url
		return connectToDB(url)
	}

	if options.port <= 0 {
		return nil, errors.New("port parameter must be more than 0")
	}
	url = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		options.user,
		options.password,
		options.host,
		options.port,
		dbName)

	if options.params != nil {
		paramAdder := addParams(url)
		if v, ok := options.params["sslmode"]; ok {
			url = paramAdder("sslmode=" + v)
		}
		if v, ok := options.params["TimeZone"]; ok {
			url = paramAdder("TimeZone=" + v)
		}
	}

	return connectToDB(url)
}

func addParams(url string) func(string) string {
	url += "?"
	return func(param string) string {
		url += param + "&"
		return url[:len(url)-1]
	}
}

func connectToDB(url string) (*database.Queries, error) {
	config := PgxPoolConfig(url)
	connPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error acquiring connnection to the database pool: ", err)
	}

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Couldn't ping database")
	}
	return database.New(connPool), nil
}
