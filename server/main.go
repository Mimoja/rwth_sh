package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const host = ""
const port = 9080

func errorResponse(c *gin.Context, code int, err string) {
	c.String(code, "error %d: %s", code, err)
}

func main() {
	hostAndPort := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting http server on %s\n", hostAndPort)

	initShortener()

	multidom := make(MultiDomainRouter)
	multidom["dashboard.localhost:9080"] = getDashboardRouter()

	defer sqliteDatabase.Close()

	if err := http.ListenAndServe(":9080", multidom); err != nil {
		log.Fatal(err)
	}
}
