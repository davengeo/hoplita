package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		key := c.DefaultQuery("key", "none")
		documentType := c.Query("document-type")

		c.String(http.StatusOK, "Looking for %s %s", key, documentType)
	})
	router.Run(":8081") // listen and server on 0.0.0.0:8080
}