package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	. "go-link-shortener/server/database"
)

type DomainRow struct {
	Domain, Short, Long, Desc string
}

func InitShortener() {
	insertURL(Database, "rwth.sh", "abc", "https://google.com", "test entry")
	insertURL(Database, "abc.rwth.sh", "abc", "https://google.com", "test entry 2")
	insertURL(Database, "o.rwth.sh", "", "https://online.rwth-aachen.de", "test entry 2")

	printStoredURLs(Database)
}

func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	url, err := getURL(Database, r.Host, r.RequestURI[1:])
	if err != nil {
		http.Error(w, "Not found", 404)
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// We are passing db reference connection from main to our method with other parameters
func insertURL(db *sql.DB, domain string, short string, long string, comment string) {
	log.Println("Inserting url record ...")
	insertStudentSQL := `INSERT INTO urls(domain, path, target, comment) VALUES (?, ?, ?, ?)`
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

func GetURLCount(db *sql.DB) int {
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

func GetStoredURLs(db *sql.DB) []DomainRow {
	row, err := db.Query("SELECT * FROM urls")
	if err != nil {
		log.Fatal(err)
	}

	defer row.Close()

	// Iterate and fetch the records from result cursor
	result := make([]DomainRow, 0)
	for row.Next() {
		var domain, short, long, desc string
		row.Scan(&domain, &short, &long, &desc)
		result = append(result, DomainRow{domain, short, long, desc})
	}
	return result
}

func printStoredURLs(db *sql.DB) {
	rows := GetStoredURLs(db)
	println("Stored URLs:", len(rows))
	for _, r := range rows {
		println(r.Domain, r.Short, " -> ", r.Long, r.Desc)
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
		return "", sql.ErrNoRows
	}
	return long, nil
}
