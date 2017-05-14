package main

import "github.com/gin-gonic/gin"

func UserMiddleware(c *gin.Context) {
	//fmt.Println("Im a dummy!")

	sid, _ := c.Cookie("sid")

	session := sessionStore.Get(sid)

	data := make(map[string]interface{})
	data["session"] = session
	c.Set("data", data)

	/*
		sessionStore.Set(session.sid, session)
		cookie := &http.Cookie{Name: "sid", Value: session.sid}

		http.SetCookie(c.Writer, cookie)
	*/

	// Pass on to the next-in-chain
	c.Next()
}
