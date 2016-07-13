package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/couchbase/gocb"
	"fmt"
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

	for {
		var document Document
		document=<-income
		go PipeLine(document)
	}
}

func PipeLine(document Document) {
	second(verify(document))
}

func verify(doc Document) <-chan Document {
	out := make(chan Document)
	myCluster, _ := gocb.Connect("couchbase://localhost")
	myBucket, _ := myCluster.OpenBucket("sync_gateway", "")
	var value interface{}
	cas, _ := myBucket.Get(doc.Id, &value)
	fmt.Printf("Got value `%+v` with CAS `%08x`\n", value, cas)
	return out
}

func second(in chan Document) <-chan Document {
	out := make(chan Document)
	<-in
	return out
}
