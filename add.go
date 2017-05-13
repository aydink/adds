package main

import (
	"log"
	"time"
)

type Add struct {
	Id          int
	Uid         int
	Title       string
	Description string
	Price       float32
	Category    int
	Region      int
	Negotiable  bool
	Created     time.Time
	Updated     time.Time
	Image       string
	ViewCount   int
}

func AddList() []Add {
	stmt, err := db.Prepare("SELECT id, uid, title, description, price, category, region, negotiable, created, updated, image, view_count FROM adds WHERE region = ?")
	if err != nil {
		log.Println(err)
	}

	rows, err := stmt.Query(1)
	if err != nil {
		log.Println(err)
	}

	adds := make([]Add, 0)

	for rows.Next() {
		add := Add{}
		err = rows.Scan(&add.Id, &add.Uid, &add.Title, &add.Description, &add.Price, &add.Category, &add.Region, &add.Negotiable, &add.Created, &add.Updated, &add.Image, &add.ViewCount)
		if err != nil {
			log.Fatal(err)
		}

		adds = append(adds, add)
	}

	return adds
}
