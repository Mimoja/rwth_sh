package dashboard

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/Masterminds/sprig"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"go-link-shortener/server/router"
)

type Page struct {
	IsAdmin  bool
	URLCount int
}

const templateFolder = "templates"

func errorResponse(c *gin.Context, code int, err string) {
	c.String(code, "error %d: %s", code, err)
}

func getAllTemplates() []string {
	templateFiles := []string{}
	files, err := ioutil.ReadDir(templateFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		templateFiles = append(templateFiles, path.Join(templateFolder, f.Name()))
	}
	fmt.Println(templateFiles)
	return templateFiles
}

var templates = template.Must(template.New("dummy").Funcs(sprig.FuncMap()).ParseFiles(getAllTemplates()...))

func display(c *gin.Context, tmpl string, data interface{}) {
	c.Header("Content-Type", "html")

	templates = template.Must(template.New("dummy").Funcs(sprig.FuncMap()).ParseFiles(getAllTemplates()...))

	err := templates.ExecuteTemplate(c.Writer, tmpl, data)
	if err != nil {
		log.Println("Template error: ", err)
		errorResponse(c, http.StatusInternalServerError, "Could not activate template")
	}
}

func adminHandler(c *gin.Context) {
	page := Page{
		IsAdmin: true,
	}

	display(c, "admin", &page)
}

func mainHandler(c *gin.Context) {
	page := Page{
		URLCount: router.GetURLCount(router.SqliteDatabase),
	}

	display(c, "main", &page)
}

func GetDashboardRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", mainHandler)

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	}))

	authorized.GET("/", adminHandler)

	r.Use(static.Serve("/static/", static.LocalFile("./static", false)))
	return r
}
