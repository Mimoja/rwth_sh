package database

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var noChangeError = errors.New("no change")

func RunMigrations(db *sql.DB) {
	log.Println("Running migrations")
	driver, _ := sqlite3.WithInstance(db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://server/database/migrations", "sqlite3", driver)

	if err != nil {
		log.Fatal(err)
	}

	previous_ver, _, _ := m.Version()

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			log.Fatal(err)
		}
	} else {
		new_ver, _, _ := m.Version()
		log.Println("migrated DB successfully from version", previous_ver, "to", new_ver)
	}
}
