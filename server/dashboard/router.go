package dashboard

import (
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-link-shortener/server/common"
	"go-link-shortener/server/dashboard/handlers"
	"go-link-shortener/server/dashboard/middleware"
	"go-link-shortener/server/globals"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", handlers.IndexGetHandler)
	g.GET("/login", handlers.LoginGetHandler)
	g.GET("/overview", handlers.PublicOverviewGetHandler)

	g.POST("/login", handlers.LoginPostHandler)
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/dashboard", handlers.IndexGetHandler)
	g.POST("/logout", handlers.LogoutGetHandler)
	g.GET("/admin", handlers.AdminGetHandler)

	g.POST("/api/url/add", handlers.AddURLPostHandler)
	g.POST("/api/url/delete", handlers.DeleteURLPostHandler)
	g.POST("/api/url/update", handlers.UpdateURLPostHandler)
}

func GetDashboardRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static/", "./static")
	// setup templating and define custom template function
	router.SetFuncMap(template.FuncMap{
		"struct2json":  common.Struct2JSON,
		"getHostname":  common.GetHostname,
		"buildAddress": common.BuildAddress,
	})
	router.LoadHTMLGlob("templates/*.tmpl")

	router.Use(sessions.Sessions("sessions", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	PrivateRoutes(private)

	return router
}
