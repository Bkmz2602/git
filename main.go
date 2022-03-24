package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"goServ/controller"
)

func main() {
	control, err := controller.NewController()
	if err != nil {
		return
	}

	defer control.GetDB().Close()

	router := gin.Default()
	router.GET("/peoples/", control.GetPeoples)
	router.GET("/peoples/:id", control.GetPeoplesById)
	router.POST("/peoples/", control.PostPeoples)
	router.PUT("/peoples/", control.ModifyPeoples)
	router.DELETE("/peoples/", control.DeletePeoplesById)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8989")
}
