package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"log"
)

type message struct {

	value interface {}

}

func main() {
	router := gin.Default()

	msg := message{}

	router.GET("/webhook", func(c *gin.Context) {
		key := c.DefaultQuery("key", "none")
		documentType := c.Query("document-type")
		c.String(http.StatusOK, "Looking for %s %s", key, documentType)
		msg.value = key
	})

	go OrdersLoop(&msg)

	router.Run(":8081")

}


func OrdersLoop(messages *message) {

	for {
		time.Sleep(300 * time.Millisecond)
		log.Printf("hello:%s", messages.value)
	}

}