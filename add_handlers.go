package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
)

func HandleAddUpdate(c *gin.Context) {
	id := c.Params.ByName("id")
	addId, _ := strconv.Atoi(id)

	data := c.MustGet("data").(map[string]interface{})
	log.Printf("%+v\n", data["session"])

	add, err := GetAddById(addId)
	if err != nil {
		log.Println(err)
	}

	data["add"] = add
	renderTemplate(c.Writer, "new_edit_add.html", data)
}

func HandlePhotoUploadForm(c *gin.Context) {
	id := c.Params.ByName("id")
	addId, _ := strconv.Atoi(id)

	data := c.MustGet("data").(map[string]interface{})
	log.Printf("%+v\n", data["session"])

	add, err := GetAddById(addId)
	if err != nil {
		log.Println(err)
	}

	data["add"] = add
	renderTemplate(c.Writer, "upload_photo.html", data)
}

// save uploaded photos for a given add
func HandlePhotoUpload(c *gin.Context) {
	data := c.MustGet("data").(map[string]interface{})
	session := data["session"].(*Session)

	log.Printf("%+v\n", session)

	strAid := c.DefaultPostForm("aid", "0")
	aid, err := strconv.Atoi(strAid)

	if err != nil {
		log.Println(err)
	}
	log.Println("add id:", aid)

	uid := session.SessionUser.Id

	if session.SessionUser.Loggedin == false {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	for i := 1; i < 5; i++ {
		file, header, err := c.Request.FormFile("file" + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		} else {
			err = saveFile(uid, aid, file, header)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func HandleAddsView(c *gin.Context) {
	id := c.Params.ByName("id")
	addId, _ := strconv.Atoi(id)

	data := c.MustGet("data").(map[string]interface{})
	log.Printf("%+v\n", data["session"])

	add, err := GetAddById(addId)
	if err != nil {
		log.Println(err)
	}

	data["photos"] = GetPhotos(addId)
	log.Println(data["photos"])
	data["add"] = add
	renderTemplate(c.Writer, "adds_view.html", data)
}

func registerAddHandlers() {

	router.GET("/adds/edit/:id", HandleAddUpdate)
	router.GET("/adds/addphoto/:id", HandlePhotoUploadForm)
	router.POST("/adds/addphoto", HandlePhotoUpload)
	router.GET("/adds/view/:id", HandleAddsView)

	router.GET("/adds/new", func(c *gin.Context) {
		data := c.MustGet("data").(map[string]interface{})
		session := data["session"].(*Session)
		log.Printf("%+v\n", session)

		if session.SessionUser.Loggedin == false {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		data["add"] = Add{}
		renderTemplate(c.Writer, "new_edit_add.html", data)
	})

	router.POST("/adds/new", func(c *gin.Context) {

		data := c.MustGet("data").(map[string]interface{})
		session := data["session"].(*Session)

		log.Printf("%+v\n", session)
		log.Println("Is User Loggedin:", session.SessionUser.Loggedin)

		if session.SessionUser.Loggedin == false {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		id := c.DefaultPostForm("id", "0")
		title := c.DefaultPostForm("title", "")
		description := c.DefaultPostForm("description", "")
		price := c.DefaultPostForm("price", "")
		negotiable := c.DefaultPostForm("negotiable", "")
		category := c.DefaultPostForm("category", "1")
		region := c.DefaultPostForm("region", "1")

		formHasError := false
		formErrors := make(map[string]string)

		idInt, err := strconv.Atoi(id)
		if err != nil {
			formHasError = true
			formErrors["id"] = "İlan numarası geçersiz."
		}

		title = strings.TrimSpace(title)
		if len(title) < 6 {
			formHasError = true
			formErrors["title"] = "İlan başlığı en az 5 karakter uzunluğunda olmalıdır."
		}

		// convert from float represantation from turkish format to golang format
		price = strings.Replace(price, ",", ".", -1)

		priceFloat, err := strconv.ParseFloat(price, 32)
		if err != nil {
			formHasError = true
			formErrors["price"] = "Geçerli bir sayı giriniz. Örnek: 12500 veya 2,5"
		}

		negotiableValue := false
		if negotiable == "on" {
			negotiableValue = true
		}

		categoryInt, err := strconv.Atoi(category)
		if err != nil {
			formHasError = true
			formErrors["category"] = "İlan için bir kategori seçin."
		}

		regionInt, err := strconv.Atoi(region)
		if err != nil {
			formHasError = true
			formErrors["region"] = "Lütfen lojman bölgesini seçin."
		}

		add := Add{}
		add.Id = idInt
		add.Uid = session.SessionUser.Id
		add.Title = title
		add.Description = description
		add.Price = float32(priceFloat)
		add.Negotiable = negotiableValue
		add.Category = categoryInt
		add.Region = regionInt

		if formHasError == false {

			// check if we are creating or editing an Add
			if idInt == 0 {
				// create a new Add
				err := CreateAdd(add)
				if err != nil {
					log.Println(err)
					formHasError = true
					formErrors["database"] = err.Error()
				}
			} else {
				// updating existing Add
				err := UpdateAdd(add)
				if err != nil {
					log.Println(err)
					formHasError = true
					formErrors["database"] = err.Error()
				}
			}

		}

		if formHasError {
			data["formErrors"] = formErrors
			log.Printf("%+v\n", formErrors)
			data["add"] = add
			renderTemplate(c.Writer, "new_edit_add.html", data)
		} else {
			c.Redirect(http.StatusFound, "/adds/myadds")
		}
	})

	router.GET("/adds/myadds", func(c *gin.Context) {

		data := c.MustGet("data").(map[string]interface{})
		session := data["session"].(*Session)

		data["adds"] = ListAddsByUser(session.SessionUser.Id)
		renderTemplate(c.Writer, "adds_myadds.html", data)
	})
}
