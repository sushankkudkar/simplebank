package main

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sushankkudkar/simplebank/api"
	db "github.com/sushankkudkar/simplebank/db/sqlc"
	"github.com/sushankkudkar/simplebank/util"
)

func main() {
	var err error
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config variable:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)
	server, _ := api.NewServer(config, &store)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
}
