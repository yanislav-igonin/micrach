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

func Migrate() {
	// m, err := migrate.New("file://migrations", Config.Db.Url)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// err = m.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	log.Panicln(err)
	// }
	// m.Close()
	// TODO: check another libraries
	log.Println("database migrations - online")
}
