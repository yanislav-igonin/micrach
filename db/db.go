package db

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	Config "micrach/config"
	Files "micrach/files"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

type MigrationsMap map[int]string
type Migration struct {
	ID   int
	Name string
}

func Init(ctx context.Context) {
	var err error
	Pool, err = pgxpool.New(ctx, Config.Db.Url)
	if err != nil {
		log.Println("database - offline")
		log.Panicln(err)
	}

	log.Println("database - online")
}

func Migrate(ctx context.Context) {
	dbMigrations := getDbMigrations(ctx)
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
			runMigration(ctx, id, name, sql)
			log.Println("migration - " + name + " - online")
		}
	}

	log.Println("migrations - online")
}

func runMigration(ctx context.Context, mid int, mname, msql string) {
	_, err := Pool.Exec(ctx, msql)
	if err != nil {
		log.Panicln(err)
	}

	_, err = Pool.Exec(ctx, `INSERT INTO migrations (id, name) VALUES ($1, $2)`, mid, mname)
	if err != nil {
		log.Panicln(err)
	}
}

func getDbMigrations(ctx context.Context) MigrationsMap {
	rows, err := Pool.Query(ctx, `SELECT id, name FROM migrations`)
	if err != nil {
		if isUndefinedTable(err) {
			return make(MigrationsMap)
		}
		log.Panicln(err)
	}
	defer rows.Close()

	migrationsMap := make(MigrationsMap)
	for rows.Next() {
		var m Migration
		err = rows.Scan(&m.ID, &m.Name)
		if err != nil {
			log.Panicln(err)
		}

		migrationsMap[m.ID] = m.Name
	}
	if err := rows.Err(); err != nil {
		log.Panicln(err)
	}

	return migrationsMap
}

func isUndefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}
