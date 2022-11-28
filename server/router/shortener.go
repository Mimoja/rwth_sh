package router

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	. "go-link-shortener/server/database"
)

var InsertUniqueError = errors.New("Unique constraint failed")

type DomainRow struct {
	Subdomain, Path, Target, Comment string
	Id                               uint32
}

func InitShortener() {
	InsertURL(Database, DomainRow{"", "abc", "https://google.com", "test entry", 0})
	InsertURL(Database, DomainRow{"abc", "abc", "https://google.com", "test entry 2", 0})
	InsertURL(Database, DomainRow{"o", "", "https://online.rwth-aachen.de", "test entry 2", 0})

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
func InsertOrUpdateURL(db *sql.DB, entry DomainRow, update bool) error {
	log.Println("Inserting url record ...")
	query := `INSERT INTO urls(subdomain, path, target, comment) VALUES (?, ?, ?, ?)`
	if update {
		query = `UPDATE urls subdomain=?, path=?, target=?, comment=? WHERE id=?`
	}

	statement, err := db.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("Prepare failed", err.Error())
	}

	if update {
		_, err = statement.Exec(entry.Subdomain, entry.Path, entry.Target, entry.Comment, entry.Id)
	} else {
		_, err = statement.Exec(entry.Subdomain, entry.Path, entry.Target, entry.Comment)
	}

	if err != nil {
		if strings.Index(err.Error(), "UNIQUE constraint failed") == 0 {
			return InsertUniqueError
		} else {
			log.Panic(err)
		}
	}
	return nil
}

func InsertURL(db *sql.DB, entry DomainRow) error {
	return InsertOrUpdateURL(db, entry, false)
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
		var domainRow DomainRow

		row.Scan(&domainRow.Subdomain, &domainRow.Path, &domainRow.Target, &domainRow.Comment, &domainRow.Id)
		result = append(result, domainRow)
	}
	return result
}

func printStoredURLs(db *sql.DB) {
	rows := GetStoredURLs(db)
	println("Stored URLs:", len(rows))
	for _, r := range rows {
		fmt.Printf("%d: %s.HOST/%s -> %s [%s]\n", r.Id, r.Subdomain, r.Path, r.Target, r.Comment)
	}
}

func getURL(db *sql.DB, domain string, short string) (string, error) {
	row := db.QueryRow("SELECT target FROM urls WHERE domain=? AND short=?", domain, short)
	if row == nil {
		return "", fmt.Errorf("Failed to querry row")
	}

	// Parse row into Activity struct
	var target string
	if err := row.Scan(target); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return "", sql.ErrNoRows
	}
	return target, nil
}
