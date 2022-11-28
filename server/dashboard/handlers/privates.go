package handlers

import (
	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"
	"go-link-shortener/server/router"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	if user == nil {
		panic("not logged in user reached this page")
	}

	log.Println("logging out user:", user)

	session.Delete(globals.Userkey)
	if err := session.Save(); err != nil {
		log.Println("Failed to delete session:", err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func AdminGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	c.HTML(http.StatusOK, "admin", gin.H{
		"URLCount": router.GetURLCount(Database),
		"user":     user,
	})
}
