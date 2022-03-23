package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type album struct {
	ID          string `json:"id"`
	Last_name   string `json:"last_name"`
	First_name  string `json:"first_name"`
	Middle_name string `json:"middle_name"`
	Address     string `json:"address"`
}

type Api interface {
	GetDB() *sql.DB
	GetAlbums(ctx *gin.Context)
	GetAlbumsById(ctx *gin.Context)
	PostAlbums(ctx *gin.Context)
	ModifyAlbums(ctx *gin.Context)
}

type controller struct {
	DB *sql.DB
}

func NewController() (Api, error) {
	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &controller{
		DB: db,
	}, nil
}

func (c *controller) GetDB() *sql.DB {
	return c.DB
}

func (c *controller) GetAlbums(ctx *gin.Context) {
	var books []album
	var book album
	sql_statement := "SELECT p.id, p.last_name, p.first_name, p.middle_name, r.address FROM people p JOIN registry r ON p.id = r.people_id;"
	rows, _ := c.DB.Query(sql_statement)

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&book.ID, &book.Last_name, &book.First_name, &book.Middle_name, &book.Address)
		books = append(books, book)
	}

	ctx.IndentedJSON(http.StatusOK, books)
}

func (c *controller) GetAlbumsById(ctx *gin.Context) {
	p := ctx.Param("id")

	var books []album
	var book album

	sql_statement := "select p.id, p.last_name, p.first_name, p.middle_name, r.address From people p JOIN registry r on p.id = r.people_id where r.people_id = $1;"
	rows, _ := c.DB.Query(sql_statement, p)

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&book.ID, &book.Last_name, &book.First_name, &book.Middle_name, &book.Address)
		books = append(books, book)
	}

	if book.ID == p {
		ctx.IndentedJSON(http.StatusOK, books)
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func (c *controller) PostAlbums(ctx *gin.Context) {
	var newAlbum album
	ctx.BindJSON(&newAlbum)

	insertPeople := "insert into people (last_name, first_name, middle_name) VALUES ($1, $2, $3);"
	c.DB.Query(insertPeople, newAlbum.Last_name, newAlbum.First_name, newAlbum.Middle_name)

	insertRegistry := "insert into registry(people_id, address) values ((select max(people.id) from people),$1);"
	c.DB.Query(insertRegistry, newAlbum.Address)

	ctx.IndentedJSON(http.StatusCreated, newAlbum)
}

func (c *controller) ModifyAlbums(ctx *gin.Context) {
	id := ctx.Param("id")
	newAddress, _ := ctx.GetQuery("address")

	modifyRegistry := "UPDATE registry r SET address = $1 WHERE people_id = $2;"
	c.DB.Query(modifyRegistry, newAddress, id)
	c.GetAlbumsById(ctx)
}
