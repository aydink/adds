package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var sessionStore *SessionStore

func main() {
	var err error
	db, err = sql.Open("mysql", "root:sanane@/adds?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := gin.Default()
	// provide request user
	router.Use(UserMiddleware)

	router.GET("/", func(c *gin.Context) {

		data := c.MustGet("data").(map[string]interface{})

		data["adds"] = AddList()
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

		saveFile(file, header)

	})

	router.GET("/login", func(c *gin.Context) {
		data := make(map[string]interface{})

		renderTemplate(c.Writer, "login.html", data)
	})

	router.POST("/login", func(c *gin.Context) {

		email := c.DefaultPostForm("email", "")
		password := c.DefaultPostForm("password", "")

		if user, ok := LoginUser(email, password); ok {
			user.Loggedin = true

			session := NewSession(user, make(map[interface{}]interface{}))
			//fmt.Println(session)
			session.isNew = false
			sessionStore.Set(session.sid, session)
			cookie := &http.Cookie{Name: "sid", Value: session.sid}

			http.SetCookie(c.Writer, cookie)

			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			data := make(map[string]interface{})
			data["errorMessage"] = "Geçersiz kullanıcı adı veya şifre"
			renderTemplate(c.Writer, "login.html", data)
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

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}

func init() {
	sessionStore = NewSessionStore()
}
