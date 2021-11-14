package db

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	Config "micrach/config"
	Files "micrach/files"

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
	migrations := Files.GetFullFilePathsInFolder("migrations")
	log.Println(migrations)
	for _, m := range migrations {
		filename := filepath.Base(m)
		splitted := strings.Split(filename, "-")
		id, name := splitted[0], splitted[1]
		log.Println(id, name)
		log.Println(Files.ReadFileText(m))
	}
	log.Println("database migrations - online")
}
