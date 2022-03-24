package controller

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type people struct {
	ID          string `json:"id"`
	Last_name   string `json:"last_name"`
	First_name  string `json:"first_name"`
	Middle_name string `json:"middle_name"`
	Address     string `json:"address"`
}

type Api interface {
	GetDB() *sql.DB
	GetPeoples(ctx *gin.Context)
	GetPeoplesById(ctx *gin.Context)
	PostPeoples(ctx *gin.Context)
	ModifyPeoples(ctx *gin.Context)
	DeletePeoplesById(ctx *gin.Context)
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

func (c *controller) GetPeoples(ctx *gin.Context) {
	var lists []people
	var list people
	page, _ := ctx.GetQuery("page")
	sortBy, _ := ctx.GetQuery("sort_by")
	limit, _ := ctx.GetQuery("limit")
	var limitDefault uint64 //limit to print per page

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("p.id", "p.last_name", "p.first_name", "p.middle_name", "r.address").
		From("people as p").
		Join("registry as r on r.people_id = p.id")

	//limits result
	if limit != "" {
		if lim, err := strconv.Atoi(limit); err == nil {
			builder = builder.Limit(uint64(lim))
			limitDefault = uint64(lim)
		}
	} else {
		limitDefault = 5
		builder = builder.Limit(limitDefault)
	}

	//pagination
	if pg, err := strconv.Atoi(page); err == nil {
		if pg > 1 {
			offs := pg * int(limitDefault)
			builder = builder.Offset(uint64(offs))
		}
	}

	//sorting
	if sortBy != "" {
		builder = builder.OrderBy(sortBy)
	} else {
		builder = builder.OrderBy("p.id")
	}

	req, _, err := builder.ToSql()
	fmt.Println(req)

	if err != nil {
		fmt.Printf("%v,sql", err)
		return
	}

	rows, err := c.DB.Query(req)

	if err != nil {
		fmt.Printf("%v,rows", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&list.ID, &list.Last_name, &list.First_name, &list.Middle_name, &list.Address)
		lists = append(lists, list)
	}

	ctx.IndentedJSON(http.StatusOK, lists)
}

func (c *controller) GetPeoplesById(ctx *gin.Context) {
	id := ctx.Param("id")

	var lists []people
	var list people

	sql_statement := "select p.id, p.last_name, p.first_name, p.middle_name, r.address From people p JOIN registry r on p.id = r.people_id where r.people_id = $1;"
	rows, _ := c.DB.Query(sql_statement, id)

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&list.ID, &list.Last_name, &list.First_name, &list.Middle_name, &list.Address)
		lists = append(lists, list)
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "people not found"})
}

func (c *controller) PostPeoples(ctx *gin.Context) {
	var newPeople people
	ctx.BindJSON(&newPeople)

	insertPeople := "insert into people (last_name, first_name, middle_name) VALUES ($1, $2, $3);"
	c.DB.Query(insertPeople, newPeople.Last_name, newPeople.First_name, newPeople.Middle_name)

	insertRegistry := "insert into registry(people_id, address) values ((select max(people.id) from people),$1);"
	c.DB.Query(insertRegistry, newPeople.Address)

	c.GetPeoples(ctx)
	//	ctx.IndentedJSON(http.StatusCreated, newPeople)
}

func (c *controller) ModifyPeoples(ctx *gin.Context) {
	var changePeopleAddress people
	ctx.BindJSON(&changePeopleAddress)

	id := changePeopleAddress.ID
	newAddress := changePeopleAddress.Address

	modifyRegistry := "UPDATE registry r SET address = $1 WHERE people_id = $2;"
	c.DB.Query(modifyRegistry, newAddress, id)
	c.GetPeoples(ctx)
}

func (c *controller) DeletePeoplesById(ctx *gin.Context) {
	var id people
	ctx.BindJSON(&id)

	deleteRequest := "DELETE FROM people WHERE id = $1;"
	c.DB.Query(deleteRequest, id.ID)
	c.GetPeoples(ctx)
}
