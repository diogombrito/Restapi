package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	idP     int
	user    string
	pass    string
	ageP    string
	nameP   string
	familyP string
	roleP   string
}

func selectID(uid int64) person {
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM persons WHERE id = ?", uid)
	checkErr(err)
	var id int
	var username string
	var password string
	var age string
	var name string
	var family string
	var role string
	var newp person
	for rows.Next() {
		err = rows.Scan(&id, &username, &age, &name, &family, &password, &role)
		//err = rows.Scan(&id, &username, &password)
		checkErr(err)
		//newp = person{idP: &id, user: username, pass: password}

		newp = person{idP: id, user: username, pass: password, ageP: age, nameP: name, familyP: family, roleP: role}
	}
	rows.Close()
	db.Close()
	return newp
}

func selectUsername(usern string) person {
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM persons WHERE username = ?", usern)
	checkErr(err)
	var id int
	var username string
	var password string
	var age string
	var name string
	var family string
	var role string
	var newp person
	for rows.Next() {
		err = rows.Scan(&id, &username, &age, &name, &family, &password, &role)
		//err = rows.Scan(&id, &username, &password)
		checkErr(err)
		//newp = person{idP: &id, user: username, pass: password}

		newp = person{idP: id, user: username, pass: password, ageP: age, nameP: name, familyP: family, roleP: role}
	}
	rows.Close()
	db.Close()
	return newp
}

func selectAll() []person {
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	rows, err := db.Query("SELECT id,username,age,name,family,role FROM persons")
	checkErr(err)
	var id int
	var username string
	var age string
	var name string
	var family string
	var role string
	allperson := []person{}
	for rows.Next() {
		err = rows.Scan(&id, &username, &age, &name, &family, &role)
		checkErr(err)
		allperson = append(allperson, person{idP: id, user: username, ageP: age, nameP: name, familyP: family, roleP: role})
	}
	rows.Close()
	db.Close()
	return allperson
}

func insertID(usename string, age string, name string, family string, password string, role string) (int64, bool) {
	// insert
	// lidar com unique username
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	stmt, err := db.Prepare("INSERT INTO persons (username, age, name,family,password,role) VALUES (?,?,?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(usename, age, name, family, password, role)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
	db.Close()
	return id, checkErr(err)
}

func deleteID(username string) bool {
	// delete
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	stmt, err := db.Prepare("DELETE from persons where username=?")
	checkErr(err)
	res, err := stmt.Exec(username)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	return true

}

func updateID(id int, username string, name string, family string, role string) bool {
	// insert
	db, err := sql.Open("sqlite3", "DataBase/foo.db")
	checkErr(err)
	stmt, err := db.Prepare("update persons set username = ?,name = ?, family = ?, role = ? where id=?")
	checkErr(err)
	res, err := stmt.Exec(username, name, family, role, id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	db.Close()
	return true
}

func testdb() {

	database, _ := sql.Open("sqlite3", "DataBase/foo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS persons (id int NOT NULL PRIMARY KEY, username text,age date,name text,family text,password text,role text)")
	statement.Exec()
}

func checkErr(err error) bool {
	if err != nil {
		panic(err)
		return false
	}
	return true
}
