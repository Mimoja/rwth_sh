package handlers

import (
	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"
	"go-link-shortener/server/router"
	"log"
	"sort"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"net/http"
)

func IndexGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	c.HTML(http.StatusOK, "main", gin.H{
		"URLCount": router.GetURLCount(Database),
		"user":     user,
	})
}

func LoginGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	if user != nil {
		log.Printf("found user %s\n", user)
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"message": "Please logout first",
		})
	} else {
		c.HTML(http.StatusOK, "login", gin.H{
			"message": "",
		})
	}
}

func PublicOverviewGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	links := router.GetStoredURLs(Database)
	// filter array for elements that are public
	pub_links := links[:0]
	for _, e := range links {
		if e.IsPublic {
			pub_links = append(pub_links, e)
		}
	}
	sort.Slice(pub_links, func(i, j int) bool {
		if pub_links[i].Subdomain != pub_links[j].Subdomain {
			return pub_links[i].Subdomain < pub_links[j].Subdomain
		}
		return pub_links[i].Path < pub_links[j].Path
	})

	c.HTML(http.StatusOK, "linkOverview", gin.H{
		"URLCount": router.GetURLCount(Database),
		"user":     user,
		"links":    pub_links,
	})
}
