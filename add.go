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

func ListAdds() []Add {
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

func ListAddsByUser(uid int) []Add {
	stmt, err := db.Prepare("SELECT id, uid, title, description, price, category, region, negotiable, created, updated, image, view_count FROM adds WHERE uid = ?")
	if err != nil {
		log.Println(err)
	}

	rows, err := stmt.Query(uid)
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

func GetAddById(id int) (Add, error) {
	stmt, err := db.Prepare("SELECT id, uid, title, description, price, category, region, negotiable, created, updated, image, view_count FROM adds WHERE id = ?")
	if err != nil {
		log.Println(err)
	}

	add := Add{}

	err = stmt.QueryRow(id).Scan(&add.Id, &add.Uid, &add.Title, &add.Description, &add.Price, &add.Category, &add.Region, &add.Negotiable, &add.Created, &add.Updated, &add.Image, &add.ViewCount)
	if err != nil {
		log.Println(err)
		return add, err
	}

	return add, nil
}

func CreateAdd(add Add) error {
	stmt, err := db.Prepare("INSERT INTO adds (uid, title, description, price, category, region, negotiable) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Exec(add.Uid, add.Title, add.Description, add.Price, add.Category, add.Region, add.Negotiable)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	add.Id = int(id)

	return nil
}

func UpdateAdd(add Add) error {
	stmt, err := db.Prepare("UPDATE adds SET title=?, description=?, price=?, category=?, region=?, negotiable=? WHERE id=? AND uid=?")
	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Exec(add.Title, add.Description, add.Price, add.Category, add.Region, add.Negotiable, add.Id, add.Uid)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected < 1 {
		log.Printf("Add update failed:%+v\n", add)
	}
	return nil
}
