package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type album struct {
	ID          string `json:"id"`
	Last_name   string `json:"last_name"`
	First_name  string `json:"first_name"`
	Middle_name string `json:"middle_name"`
	Address     string `json:"address"`
}

//func postAlbums(c *gin.Context) {
//	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
//	db, _ := sql.Open("postgres", connStr)
//	defer db.Close()
//
//	var newAlbum album
//
//	if err := c.BindJSON(&newAlbum); err != nil {
//		return
//	}
//
//	//albums = append(albums, newAlbum)
//	c.IndentedJSON(http.StatusCreated, newAlbum)
//	fmt.Println(newAlbum)
//
//	add(db, newAlbum.Last_name, newAlbum.First_name, newAlbum.Middle_name, newAlbum.Address)
//
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

//func getAlbumByID(c *gin.Context) {
//	id := c.Param("id")
//
//	for _, a := range albums {
//		if a.ID == id {
//			c.IndentedJSON(http.StatusOK, a)
//			return
//		}
//	}
//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
//}

/////////////////////////////////////////////////////////////////////////

func show() []album {

	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	checkError(err)

	sql_statement := "SELECT p.id, p.last_name, p.first_name, p.middle_name, r.address FROM people p JOIN registry r ON p.id = r.people_id;"
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()

	var books []album
	var book album

	for rows.Next() {
		rows.Scan(&book.ID, &book.Last_name, &book.First_name, &book.Middle_name, &book.Address)
		books = append(books, book)
	}
	return books
}

//getAlbums
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, show())
}

func main() {

	//connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	//db, err := sql.Open("postgres", connStr)
	//defer db.Close()
	//checkError(err)

	//показать всех
	//SelectResults := show(db)
	//fmt.Println(SelectResults)

	//WORK внести данные в таблицу
	//lname := "Sergeev"
	//fname := "Sergey"
	//mname := "Sergeevich"
	//adress := "St.Peterburg"
	//add(db, lname, fname, mname, adress)

	router := gin.Default()
	router.GET("/albums", getAlbums)
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/albums", postAlbums)

	router.Run("localhost:8989")
}
