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
	router.GET("/albums/", control.GetPeoples)
	router.GET("/albums/", control.GetPeoplesById)
	router.POST("/albums/", control.PostPeoples)
	router.PUT("/albums/", control.ModifyPeoples)
	router.DELETE("/albums/", control.DeletePeoplesById)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8989")
}
