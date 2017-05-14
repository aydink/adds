package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAuthHandlers() {

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

			c.Redirect(http.StatusFound, "/")
			return
		} else {
			data := make(map[string]interface{})
			data["errorMessage"] = "Geçersiz kullanıcı adı veya şifre"
			renderTemplate(c.Writer, "login.html", data)
		}
	})
}
