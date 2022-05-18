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
		// Get name without extension
		name := strings.Split(splitted[1], ".")[0]

		_, isMigrationInDb := dbMigrations[id]
		if !isMigrationInDb {
			sql := Files.ReadFileText(m)
			runMigration(id, name, sql)
			log.Println("migration - " + name + " - online")
		}
	}

	log.Println("migrations - online")
}

func runMigration(mid int, mname, msql string) {
	_, err := Pool.Exec(context.TODO(), msql)
	if err != nil {
		log.Panicln(err)
	}

	sql := `INSERT INTO migrations (id, name) VALUES ($1, $2)`
	_, err = Pool.Query(context.TODO(), sql, mid, mname)
	if err != nil {
		log.Panicln(err)
	}
}

func getDbMigrations() MigrationsMap {
	sql := `SELECT id, name FROM migrations`
	rows, err := Pool.Query(context.TODO(), sql)
	if err != nil && isNotNonExistentMigrationsTable(err) {
		log.Panicln(err)
	}

	if rows.Err() != nil && isNotNonExistentMigrationsTable(rows.Err()) {
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

func isNotNonExistentMigrationsTable(err error) bool {
	return err.Error() != `ERROR: relation "migrations" does not exist (SQLSTATE 42P01)`
}
