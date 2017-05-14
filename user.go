package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Username string
	Email    string
	Password []byte
	Token    string
	Disabled bool
	Admin    bool
	Loggedin bool
}

// GetUser Girilen kullanıcı adı şifresine uygun bir kullanıcı varsa o kulanıcıyı,
// yok ise boş bir kullanıcı ve hata döndürür.
func LoginUser(email, password string) (User, bool) {

	stmt, err := db.Prepare("SELECT id, username, password, email, token FROM users WHERE email=?")
	if err != nil {
		fmt.Println(err)
	}

	row := stmt.QueryRow(email)

	user := User{}

	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Token)
	if err != nil {
		fmt.Println(err)
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(password)) == nil {
		user.Password = []byte("")
		return user, true
	} else {
		user.Password = []byte("")
		return user, false
	}
}

// GetUser Girilen kullanıcı adı şifresine uygun bir kullanıcı varsa o kulanıcıyı,
// yok ise boş bir kullanıcı ve hata döndürür.
func CreateUser(user User) error {

	hash, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)

	stmt, err := db.Prepare("INSERT INTO users (username, password, email, token) VALUES (?, ?, ? ,?)")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(user.Username, hash, user.Email, user.Token)
	if err != nil {
		fmt.Println(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = int(id)

	return nil
}

func CreateSampleUsers() {
	user1 := User{}
	user1.Email = "aydinkilic@gmail.com"
	user1.Username = "Aydın KILIÇ"
	user1.Password = []byte("sanane")
	user1.Token = "token1"

	user2 := User{}
	user2.Email = "user2@gmail.com"
	user2.Username = "Ahmet MEHMET"
	user2.Password = []byte("banane")
	user2.Token = "token2"

	CreateUser(user1)
	CreateUser(user2)

}
