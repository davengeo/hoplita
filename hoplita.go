package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	router.GET("/webhook", func(c *gin.Context) {
		key := c.DefaultQuery("key", "none")
		documentType := c.Query("document-type")
		c.String(http.StatusOK, "Looking for %s %s", key, documentType)
		//lock.Lock()
		msg.value = key
		income<-msg
		//lock.Unlock()
	})

	go OrdersLoop(income)

	router.Run(":8081")

}


func OrdersLoop(income chan message) {

	messages := message{}
	for {
		time.Sleep(300 * time.Millisecond)
		if messages.value != nil {
			log.Printf("hello:%s", messages.value)
		}
		messages=<-income
	}

}

