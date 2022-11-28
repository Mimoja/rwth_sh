package database

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var Database *sql.DB

func InitDatabase(path string) {
	log.Printf("opening (or creating) %s in %s\n", filepath.Base(path), filepath.Dir(path))
	database, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("database created")

	RunMigrations(database)

	Database = database
}
