package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)


type Document struct {
	Id  string `json:"_id" binding:"required"`
	Rev string `json:"_rev" binding:"required"`
	Title string `json:"title"`
}

func main() {
	income := make(chan Document)

	go EventLoop(income)

	GinEngine(income).Run(":8081")
}

func GinEngine(income chan Document) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/webhook", func(c *gin.Context) {

		var dataModel Document

		if c.BindJSON(&dataModel) == nil {
			c.JSON(http.StatusAccepted, gin.H{})
			income<-dataModel
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

	})

	return router
}




func EventLoop(income chan Document) {

	message := Document{}
	i := int16(0)
	for {
		message =<-income
		i++
		log.Printf("%d %s\n", i, message.Id)
	}

}

