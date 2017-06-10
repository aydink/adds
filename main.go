package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var sessionStore *SessionStore
var router *gin.Engine

func main() {
	var err error
	db, err = sql.Open("mysql", "root:sanane@/adds?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router = gin.Default()
	// provide request user
	router.Use(UserMiddleware)

	router.GET("/", func(c *gin.Context) {

		data := c.MustGet("data").(map[string]interface{})

		data["adds"] = ListAdds()
		log.Printf("session:%+v", data["session"])

		renderTemplate(c.Writer, "index.html", data)
	})

	router.GET("/upload", func(c *gin.Context) {
		fmt.Println("upload")
		data := make(map[string]interface{})

		renderTemplate(c.Writer, "upload.html", data)
	})

	router.POST("/upload", func(c *gin.Context) {

		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			fmt.Println(err)
		}

		err = saveFile(1, 1, file, header)
		if err != nil {
			fmt.Println(err)
		}

	})

	router.GET("/createusers", func(c *gin.Context) {
		CreateSampleUsers()
		c.String(200, "Creates test users")
	})

	router.GET("/test", func(c *gin.Context) {
		data := c.MustGet("data").(map[string]interface{})
		c.String(200, "%+v", data["session"])
	})

	router.Static("/static", "./static")
	//router.StaticFS("/more_static", http.Dir("my_file_system"))
	//router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// put Add related handler to add_handlers.go file
	registerAddHandlers()
	// register login/logout handlers
	registerAuthHandlers()

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}

func init() {
	sessionStore = NewSessionStore()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
