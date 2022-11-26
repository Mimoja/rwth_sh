package dashboard

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-link-shortener/server/dashboard/handlers"
	"go-link-shortener/server/dashboard/middleware"
	"go-link-shortener/server/globals"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", handlers.IndexGetHandler)
	g.GET("/login", handlers.LoginGetHandler)
	g.POST("/login", handlers.LoginPostHandler)
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/dashboard", handlers.IndexGetHandler)
	g.GET("/logout", handlers.LogoutGetHandler)
}

func GetDashboardRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Static("/static/", "./static")
	router.LoadHTMLGlob("templates/*.tmpl")

	router.Use(sessions.Sessions("sessions", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	PrivateRoutes(private)

	return router
}
