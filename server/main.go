package main

import (
	. "go-link-shortener/server/database"

	"fmt"
	"go-link-shortener/server/dashboard"
	. "go-link-shortener/server/globals"
	. "go-link-shortener/server/router"
	"log"
	"net/http"
)

func main() {
	appConf := ConfigInit("config.yaml")

	hostAndPort := fmt.Sprintf("%s:%d", appConf.Server.Hostname, appConf.Server.Port)
	log.Printf("Starting http server on %s\n", hostAndPort)

	InitDatabase()
	InitShortener()

	multidom := make(MultiDomainRouter)

	dashboard_url := fmt.Sprintf("%s.%s", appConf.Dashboard.Subdomain, hostAndPort)
	multidom[dashboard_url] = dashboard.GetDashboardRouter()

	defer Database.Close()

	sPort := fmt.Sprintf(":%d", appConf.Server.Port)
	if err := http.ListenAndServe(sPort, multidom); err != nil {
		log.Fatal(err)
	}
}
