package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type message struct {
	key string
	value interface {}

}

func main() {
	router := gin.Default()

	msg := message{}
	//var lock sync.Mutex

	income := make(chan message)
	go OrdersLoop(income)

	router.GET("/webhook", func(c *gin.Context) {
		key := c.DefaultQuery("key", "none")
		documentType := c.Query("document-type")
		c.String(http.StatusOK, "Looking for %s %s", key, documentType)
		msg.key = key
		msg.value = documentType
		income<-msg
	})

	router.Run(":8081")
}


func OrdersLoop(income chan message) {

	messages := message{}
	for {
		if messages.value != nil {
			log.Printf("hello:%s", messages.key)
		}
		messages=<-income
	}

}

