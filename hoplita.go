package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type message struct {
	key string
	documentType string
	value interface {}

}

func main() {
	router := gin.Default()

	msg := message{}

	income := make(chan message)
	go OrdersLoop(income)

	router.GET("/webhook", func(c *gin.Context) {
		key := c.DefaultQuery("key", "none")
		documentType := c.Query("document-type")
		c.String(http.StatusOK, "Looking for %s %s", key, documentType)
		msg.key = key
		msg.documentType = documentType
		income<-msg
	})

	router.Run(":8081")
}


func OrdersLoop(income chan message) {

	messages := message{}
	i := int16(0)
	for {
		if messages.key != "" {
			i++
			log.Printf("%d\n", i)
		}
		messages=<-income
	}

}

