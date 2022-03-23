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
	router.GET("/albums/", control.GetAlbums)
	router.GET("/albums/:id", control.GetAlbumsById)
	router.POST("/albums/", control.PostAlbums)
	router.PUT("/albums/:id", control.ModifyAlbums)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8989")
}
