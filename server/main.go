package main

import (
	"fmt"
	"log"
	"net/http"

	"go-link-shortener/server/dashboard"
	. "go-link-shortener/server/router"
)

const host = ""
const port = 9080

func main() {
	hostAndPort := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting http server on %s\n", hostAndPort)

	InitShortener()

	multidom := make(MultiDomainRouter)
	multidom["dashboard.localhost:9080"] = dashboard.GetDashboardRouter()

	defer SqliteDatabase.Close()

	if err := http.ListenAndServe(":9080", multidom); err != nil {
		log.Fatal(err)
	}
}
