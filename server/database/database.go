package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var Database *sql.DB

func InitDatabase() {
	log.Println("creating sqlite-database.db")
	database, err := sql.Open("sqlite3", "./sqlite-database.db")

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("sqlite-database.db created")
	Database = database
}
