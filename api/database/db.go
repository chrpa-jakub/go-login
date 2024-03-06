package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var DB *redis.Client

func Run() {
  conn, err := redis.ParseURL(os.Getenv("DB_CONN"))

	if err != nil {
		panic("Could not parse connection string.")
	}

	DB = redis.NewClient(conn)

	if DB.Ping(context.Background()).Err() != nil {
		panic("Could not connect to DB!")
	}

}
