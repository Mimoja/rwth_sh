package main

import (
	"go-link-shortener/server/common"
	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"

	"fmt"
	"go-link-shortener/server/dashboard"
	. "go-link-shortener/server/router"
	"log"
	"net/http"
)

func main() {
	appConf := globals.ConfigInit("config.yaml")
	globals.Config = *appConf

	log.Printf("Starting http server on %s\n", common.GetHostname())

	InitDatabase(appConf.Database.Path)
	InitShortener()

	multidom := make(MultiDomainRouter)

	dashboard_url := common.GetHostname()
	if appConf.Dashboard.Subdomain != "" {
		dashboard_url = fmt.Sprintf("%s.%s", appConf.Dashboard.Subdomain, dashboard_url)
	}
	multidom[dashboard_url] = dashboard.GetDashboardRouter()

	defer Database.Close()

	sPort := fmt.Sprintf(":%d", appConf.Server.Port)
	if err := http.ListenAndServe(sPort, multidom); err != nil {
		log.Fatal(err)
	}
}
