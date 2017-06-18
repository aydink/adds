package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type Messege struct {
	AddId int    `json:"aid,omitempty"`
	Text  string `json:"text,omitempty"`
}

func GetAddOwner(aid int) (int, error) {

	stmt, err := db.Prepare("SELECT uid FROM adds WHERE id = ?")
	if err != nil {
		log.Println(err)
	}

	var uid int

	err = stmt.QueryRow(aid).Scan(&uid)
	if err != nil {
		log.Println(err)
		return uid, err
	}

	return uid, nil
}

func HandleMessage(c *gin.Context) {

	data := c.MustGet("data").(map[string]interface{})
	session := data["session"].(*Session)

	message := Messege{}

	decoder := json.NewDecoder(c.Request.Body)

	err := decoder.Decode(&message)

	if err != nil {
		log.Println(err)
	}

	log.Println(message)

	// read user id for sender
	fromUid := session.SessionUser.Id

	// find user id for the add
	toUid, err := GetAddOwner(message.AddId)
	if err != nil {
		log.Println(err)
	}

	log.Println("from:", fromUid, "to:", toUid, "add:", message.AddId, "message:", message.Text)

	response := "ok"

	if fromUid != 0 {
		err = saveMessage(message.AddId, fromUid, toUid, message.Text)
		if err != nil {
			log.Println(err)
			response = "Mesaj gönderilemedi, lütfen tekrar deneyin." + err.Error()
		}
	} else {
		response = "Mesaj gönderebilmek için lütfen giriş yapın."
	}

	c.JSON(200, gin.H{
		"result": response,
	})
}

func saveMessage(aid, fromUid, toUid int, text string) error {
	stmt, err := db.Prepare("INSERT INTO messages (aid, uid_from, uid_to, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(aid, fromUid, toUid, text)
	if err != nil {
		return err
	}

	return nil
}
