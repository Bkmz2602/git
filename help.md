package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

//type record struct {
//	lname, fname, mname, address string
//}

//func checkError(err error) {
//	if err != nil {
//		panic(err)
//	}
//}

//func show(db *sql.DB) []record {
//	sql_statement := "SELECT p.last_name, p.first_name, p.middle_name, r.address FROM people p JOIN registry r ON p.id = r.people_id;"
//	rows, err := db.Query(sql_statement)
//	checkError(err)
//	defer rows.Close()
//
//	var books []record
//	var book record
//
//	for rows.Next() {
//		rows.Scan(&book.lname, &book.fname, &book.mname, &book.address)
//		books = append(books, book)
//	}
//	return books
//}

//func add(db *sql.DB, ln string, fn string, mn string, ad string) {
//
//	sqlStatement := "INSERT INTO people (last_name, first_name, middle_name) VALUES ($1, $2, $3)"
//	db.Exec(sqlStatement, ln, fn, mn)
//
//	sqlStatement2 := "insert into registry (people_id, adress) values ((select max(people.id) from people),$1);"
//	db.Exec(sqlStatement2, ad)
//
//}

func del(db *sql.DB, ln string, fn string, mn string) {
	sqlStatements := "DELETE FROM people WHERE first_name = $1 AND last_name = $2 AND middle_name = $3"
	db.Exec(sqlStatements, ln, fn, mn)
}

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func main2() {

	//connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	//db, err := sql.Open("postgres", connStr)
	//defer db.Close()
	//checkError(err)
	//
	////WORK запросить из таблицы id, ФИО и адрес
	//SelectResults := show(db)
	//fmt.Println(SelectResults)

	//WORK внести данные в таблицу
	lname := "Sergeev"
	fname := "Sergey"
	mname := "Sergeevich"
	adress := "St.Peterburg"
	add(db, lname, fname, mname, adress)

	//WORK удалить человека из базы
	fname_toDelete := "Sergey"
	lname_toDelete := "Sergeev"
	mname_toDelete := "Sergeevich"
	del(db, fname_toDelete, lname_toDelete, mname_toDelete)

	//посмотреть что получилось
	SelectResults := show(db)
	fmt.Println(SelectResults)

	fmt.Println("Done")
}
