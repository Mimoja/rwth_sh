package handlers

import (
	. "go-link-shortener/server/database"
	"go-link-shortener/server/globals"
	"go-link-shortener/server/router"

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

}
