package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const host = ""
const port = 9080

func errorResponse(c *gin.Context, code int, err string) {
	c.String(code, "error %d: %s", code, err)
}

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

func main() {
	hostAndPort := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting http server on %s\n", hostAndPort)

	r := gin.Default()
	r.GET("/", mainHandler)
	r.Use(static.Serve("/static/", static.LocalFile("./static", false)))

	initShortener()

	defer sqliteDatabase.Close()

	multidom := make(MultiDomainRouter)
	multidom["dashboard.localhost:9080"] = r

	if err := http.ListenAndServe(":9080", multidom); err != nil {
		log.Fatal(err)
	}
}
