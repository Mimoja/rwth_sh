package handlers

import (
	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"
	"go-link-shortener/server/router"
	"log"

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

func LoginPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user != nil {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"message": "Please logout first",
		})
		return
	}

	username := c.PostForm("username")
	// password := c.PostForm("password")

	log.Println("logging in user:", username, c.Request.PostForm)

	session.Set(globals.Userkey, username)

	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login", gin.H{
			"message": "Failed to save session",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/dashboard")
}
