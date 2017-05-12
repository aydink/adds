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

func saveFile(file multipart.File, header *multipart.FileHeader) {
	filename := header.Filename
	//fmt.Println(header.Filename)

	out, err := os.Create("./tmp/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	h := md5.New()

	_, err = io.Copy(out, io.TeeReader(file, h))
	if err != nil {
		log.Fatal(err)
	}

	destinationFileName := hex.EncodeToString(h.Sum(nil))
	destinationFolder := "static/images/" + string(destinationFileName[0:2])
	os.Mkdir(destinationFolder, 0777)

	img, err := imaging.Open("./tmp/" + filename)
	if err != nil {
		panic(err)
	}
	thumb := imaging.Fill(img, 250, 250, imaging.Center, imaging.Box)

	imaging.Save(thumb, destinationFolder+"/"+destinationFileName+".jpg")

}
