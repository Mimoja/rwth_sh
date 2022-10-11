package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Page struct {
	URLCount int
}

const templateFolder = "templates"

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

func mainHandler(c *gin.Context) {
	page := Page{
		URLCount: getURLCount(sqliteDatabase),
	}

	display(c, "main", &page)
}

func getDashboardRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", mainHandler)
	r.Use(static.Serve("/static/", static.LocalFile("./static", false)))
	return r
}
