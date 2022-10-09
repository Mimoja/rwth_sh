package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
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

	insertURL(sqliteDatabase, "abc", "https://google.com", "test entry")

	getURLs(sqliteDatabase)
}

func shortenerIDHandler(c *gin.Context) {
	query := c.Params.ByName("shortenerID")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Link not specified")
		return
	}

	url, err := getURL(sqliteDatabase, query)
	if err != nil {
		errorResponse(c, 404, "Not found")
	}

	c.Redirect(http.StatusFound, url)
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE urls (
		"short" TEXT NOT NULL PRIMARY KEY,
		"long" TEXT,
		"comment" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create url table...")
	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		return
	}
	statement.Exec() // Execute SQL Statements
	log.Println("url table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertURL(db *sql.DB, short string, long string, comment string) {
	log.Println("Inserting url record ...")
	insertStudentSQL := `INSERT INTO urls(short, long, comment) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("Prepare failed", err.Error())
	}
	_, err = statement.Exec(short, long, comment)
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

func getURL(db *sql.DB, short string) (string, error) {
	row := db.QueryRow("SELECT * FROM urls WHERE short=?", short)
	if row == nil {
		return "", fmt.Errorf("Failed to querry row")
	}

	// Parse row into Activity struct
	var long string
	var desc string
	var err error
	if err = row.Scan(&short, &long, &desc); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return "", sqlite3.ErrNotFound
	}
	return long, nil
}
