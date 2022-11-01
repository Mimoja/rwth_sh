package main

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"log"
	"net/http"
)

var sqliteDatabase *sql.DB

func initShortener() {

	log.Println("sqlite-database.db created")

	sqliteDatabase, _ = sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	createTable(sqliteDatabase)

	insertURL(sqliteDatabase, "rwth.sh", "abc", "https://google.com", "test entry")
	insertURL(sqliteDatabase, "abc.rwth.sh", "abc", "https://google.com", "test entry 2")
	insertURL(sqliteDatabase, "o.rwth.sh", "", "https://online.rwth-aachen.de", "test entry 2")

	getURLs(sqliteDatabase)
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	url, err := getURL(sqliteDatabase, r.Host, r.RequestURI[1:])
	if err != nil {
		http.Error(w, "Not found", 404)
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE urls (
		"domain" TEXT NOT NULL,
		"short" TEXT NOT NULL,
		"long" TEXT,
		"comment" TEXT,
		PRIMARY KEY ("domain", "short")
	  );` // SQL Statement for Create Table

	log.Println("Create url table...")
	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Println("error: ", err)
		return
	}
	statement.Exec() // Execute SQL Statements
	log.Println("url table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertURL(db *sql.DB, domain string, short string, long string, comment string) {
	log.Println("Inserting url record ...")
	insertStudentSQL := `INSERT INTO urls(domain, short, long, comment) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("Prepare failed", err.Error())
	}
	_, err = statement.Exec(domain, short, long, comment)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func getURLCount(db *sql.DB) int {
	row, err := db.Query("SELECT count(*) FROM urls")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var count int
		row.Scan(&count)
		return count
	}
	return -1
}

func getURLs(db *sql.DB) {
	row, err := db.Query("SELECT * FROM urls")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var short string
		var long string
		var desc string

		row.Scan(&short, &long, &desc)
		println(short, " -> ", long, desc)
	}
}

func getURL(db *sql.DB, domain string, short string) (string, error) {
	row := db.QueryRow("SELECT * FROM urls WHERE domain=? AND short=?", domain, short)
	if row == nil {
		return "", fmt.Errorf("Failed to querry row")
	}

	// Parse row into Activity struct
	var long string
	var desc string
	var err error
	if err = row.Scan(&domain, &short, &long, &desc); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return "", sqlite3.ErrNotFound
	}
	return long, nil
}