package db

import (
	"context"
	"log"

	Config "micrach/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	var err error
	Pool, err = pgxpool.Connect(context.TODO(), Config.Db.Url)
	if err != nil {
		log.Println("database - offline")
		log.Panicln(err)
	}

	log.Println("database - online")
}
