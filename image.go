package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
)

func saveFile(uid int, aid int, file multipart.File, header *multipart.FileHeader) error {
	filename := header.Filename
	//fmt.Println(header.Filename)

	out, err := os.Create("./tmp/" + filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	h := md5.New()

	_, err = io.Copy(out, io.TeeReader(file, h))
	if err != nil {
		log.Println(err)
		return err
	}

	destinationFileName := hex.EncodeToString(h.Sum(nil))
	imageFolder := "static/images/" + string(destinationFileName[0:2])
	thumbFolder := "static/thumb/" + string(destinationFileName[0:2])
	os.Mkdir(imageFolder, 0777)
	os.Mkdir(thumbFolder, 0777)

	img, err := imaging.Open("./tmp/" + filename)
	if err != nil {
		log.Println(err)
		return err
	}

	resized := imaging.Fill(img, 720, 480, imaging.Center, imaging.Box)
	thumb := imaging.Fill(img, 250, 250, imaging.Center, imaging.Box)

	err = imaging.Save(resized, imageFolder+"/"+destinationFileName+".jpg")
	if err != nil {
		log.Println(err)
		return err
	}
	err = imaging.Save(thumb, thumbFolder+"/"+destinationFileName+".jpg")
	if err != nil {
		log.Println(err)
		return err
	}

	AddPhoto(uid, aid, destinationFileName)
	SetScreenPhoto(uid, aid, destinationFileName)

	return nil
}

func AddPhoto(uid int, aid int, filename string) error {

	stmt, err := db.Prepare("INSERT INTO photos (uid, aid, filename) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(uid, aid, filename)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func SetScreenPhoto(uid int, aid int, filename string) error {

	// firt set screen=0 for every photo, none is default
	stmt, err := db.Prepare("UPDATE photos SET screen = 0 WHERE uid=? AND aid=?")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(uid, aid)
	if err != nil {
		log.Println(err)
		return err
	}

	// set default photo for
	stmt, err = db.Prepare("UPDATE photos SET screen = 1 WHERE uid=? AND aid=? AND filename = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(uid, aid, filename)
	if err != nil {
		log.Println(err)
		return err
	}

	// update image field on adds table
	// we use this image field to display defaut image for the add
	// without the need for a table join operation
	// firt set screen=0 for every photo, none is default
	stmt, err = db.Prepare("UPDATE adds SET image = ? WHERE uid=? AND id=?")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(filename, uid, aid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetPhotos(aid int) []string {
	stmt, err := db.Prepare("SELECT filename FROM photos WHERE aid = ?")
	if err != nil {
		log.Println(err)
	}

	rows, err := stmt.Query(aid)
	if err != nil {
		log.Println(err)
	}

	photos := make([]string, 0)

	for rows.Next() {
		var filename string

		err = rows.Scan(&filename)
		if err != nil {
			log.Fatal(err)
		}

		photos = append(photos, filename)
	}

	return photos
}
