package handlers

import (
	"errors"
	"go-link-shortener/server/database"
	"go-link-shortener/server/router"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type URLEntry struct {
	Id        uint32 `json:"id"`
	Subdomain string `json:"subdomain"`
	Path      string `json:"path"`
	Target    string `json:"target-url"`
	Comment   string `json:"comment"`
	IsPublic  bool   `json:"is-public"`
}

func DeleteURLPostHandler(c *gin.Context) {
	type URLDeleteEntry struct {
		ID uint32 `json:"id"`
	}

	var urlEntry URLDeleteEntry

	if err := c.BindJSON(&urlEntry); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "malformed json"})
		return
	}

	if err := router.DeleteURL(database.Database, urlEntry.ID); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func AddURLPostHandler(c *gin.Context) {
	handleURLAddOrUpdatePost(c, false)
}

func UpdateURLPostHandler(c *gin.Context) {
	handleURLAddOrUpdatePost(c, true)
}

func handleURLAddOrUpdatePost(c *gin.Context, update bool) {
	var urlEntry URLEntry

	if err := c.BindJSON(&urlEntry); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// verify data
	if err := verifyAddOrUpdateRequest(&urlEntry); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	domainRow := router.DomainRow{
		Subdomain: urlEntry.Subdomain,
		Path:      urlEntry.Path,
		Target:    urlEntry.Target,
		Comment:   urlEntry.Comment,
		IsPublic:  urlEntry.IsPublic,
		Id:        urlEntry.Id,
	}

	if err := router.InsertOrUpdateURL(database.Database, &domainRow, update); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func verifyAddOrUpdateRequest(urlEntry *URLEntry) error {
	if urlEntry.Subdomain == "" && urlEntry.Path == "" {
		return errors.New("Either Subdomain or Path must not be empty")
	}

	if _, err := url.ParseRequestURI(urlEntry.Target); err != nil {
		return errors.New("Target is invalid |" + err.Error())
	}

	return nil
}
