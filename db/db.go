package db

import (
	"context"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	Config "micrach/config"
	Files "micrach/files"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

type MigrationsMap map[int]string
type Migration struct {
	ID   int
	Name string
}

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
	dbMigrations := getDbMigrations()
	sqlMigrations := Files.GetFullFilePathsInFolder("migrations")
	for _, m := range sqlMigrations {
		filename := filepath.Base(m)
		splitted := strings.Split(filename, "-")
		id, err := strconv.Atoi(splitted[0])
		if err != nil {
			log.Panicln(err)
		}
		name := strings.Split(splitted[1], ".")[0]

		if _, ok := dbMigrations[id]; !ok {
			_, err = Pool.Query(context.TODO(), Files.ReadFileText(m))
			if err != nil {
				log.Panicln(err)
			}

			sql := `INSERT INTO migrations (id, name) VALUES ($1, $2)`
			_, err = Pool.Query(context.TODO(), sql, id, name)
			if err != nil {
				log.Panicln(err)
			}
		}
	}

	log.Println("database migrations - online")
}

func getDbMigrations() MigrationsMap {
	sql := `SELECT id, name FROM migrations`
	rows, err := Pool.Query(context.TODO(), sql)
	if err != nil {
		log.Panicln(err)
	}

	if rows.Err() != nil {
		log.Panicln(rows.Err())
	}

	migrationsMap := make(MigrationsMap)
	for rows.Next() {
		var m Migration
		err = rows.Scan(&m.ID, &m.Name)
		if err != nil {
			log.Panicln(err)
		}

		migrationsMap[m.ID] = m.Name
	}

	return migrationsMap
}
