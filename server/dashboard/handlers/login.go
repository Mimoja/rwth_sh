package handlers

import (
	"go-link-shortener/server/common"
	"go-link-shortener/server/globals"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	password := c.PostForm("password")

	if EmptyUserPass(username, password) {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"message": "Please enter a username or password",
		})
		return
	}

	if !CheckUserPass(username, password) {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"message": "Wrong username or password",
		})
		return
	}

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

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func CheckUserPass(username, password string) bool {
	if username != globals.Config.Dashboard.Admin.Username {
		return false
	} else {
		return common.CheckPasswordHash(password, globals.Config.Dashboard.Admin.Password)
	}
}
