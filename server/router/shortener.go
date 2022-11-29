package router

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"
)

var InsertUniqueError = errors.New("Unique constraint failed")
var NothingChangedError = errors.New("Nothing changed - invalid ID?")

type DomainRow struct {
	Subdomain, Path, Target, Comment string
	Id                               uint32
}

func InitShortener() {
	printStoredURLs(Database)
}

func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	idx := strings.LastIndex(r.Host, globals.Config.Server.Hostname)

	if idx == -1 {
		http.Error(w, "Request with invalid hostname", 500)
		return
	}

	subdomain := strings.TrimRight(r.Host[:idx], ".")

	if subdomain == "" && r.RequestURI[1:] == "" {
		http.Redirect(w, r, "https://"+globals.DashboardURL, http.StatusFound)
	}

	url, err := getURL(Database, subdomain, r.RequestURI[1:])
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// We are passing db reference connection from main to our method with other parameters
func InsertOrUpdateURL(db *sql.DB, entry *DomainRow, update bool) error {
	log.Println("Inserting url record ...")
	query := `INSERT INTO urls(subdomain, path, target, comment) VALUES (?, ?, ?, ?)`
	if update {
		query = `UPDATE urls SET subdomain=?, path=?, target=?, comment=? WHERE id=?`
	}

	statement, err := db.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("Prepare failed", err.Error())
	}

	var result sql.Result = nil
	if update {
		result, err = statement.Exec(
			strings.ToLower(entry.Subdomain), strings.ToLower(entry.Path),
			entry.Target, entry.Comment, entry.Id)
	} else {
		result, err = statement.Exec(
			strings.ToLower(entry.Subdomain), strings.ToLower(entry.Path),
			entry.Target, entry.Comment)
	}

	if err != nil {
		if strings.Index(err.Error(), "UNIQUE constraint failed") == 0 {
			return InsertUniqueError
		} else {
			log.Panic(err)
		}
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return NothingChangedError
	}

	return nil
}

func InsertURL(db *sql.DB, entry *DomainRow) error {
	return InsertOrUpdateURL(db, entry, false)
}

func DeleteURL(db *sql.DB, id uint32) error {
	log.Println("Deleting URL with id:", id)
	query := "DELETE FROM urls WHERE id=?"

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln("Prepare failed", err.Error())
	}

	result, err := statement.Exec(id)
	if err != nil {
		return nil
	} else if res, _ := result.RowsAffected(); res == 0 {
		return NothingChangedError
	} else {
		return nil
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

func getURL(db *sql.DB, subdomain string, path string) (string, error) {
	row := db.QueryRow("SELECT target FROM urls WHERE subdomain=? AND path=?",
		strings.ToLower(subdomain), strings.ToLower(path))
	if row == nil {
		return "", fmt.Errorf("Failed to querry row")
	}

	// Parse row into Activity struct
	var target string
	if err := row.Scan(&target); err == sql.ErrNoRows || target == "" {
		log.Printf("Id not found")
		return "", sql.ErrNoRows
	}
	return target, nil
}
