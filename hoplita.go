package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)


type DataModel struct {
	Id     string `json:"_id" binding:"required"`
	Rev string `json:"_rev" binding:"required"`
}

func main() {
	router := gin.Default()


	income := make(chan DataModel)
	go OrdersLoop(income)

	router.POST("/webhook", func(c *gin.Context) {

		var dataModel DataModel

		if c.Bind(&dataModel) == nil {
			c.String(http.StatusOK, "Looking for %s %s", dataModel.Id, dataModel.Rev)
			income<-dataModel
		}


	})

	router.Run(":8081")
}




func OrdersLoop(income chan DataModel) {

	messages := DataModel{}
	i := int16(0)
	for {
		if messages.Id != "" {
			i++
			log.Printf("%d %s\n", i, messages.Id)
		}
		messages=<-income
	}

}

