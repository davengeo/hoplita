package main

import (
	"testing"
	"gopkg.in/appleboy/gofight.v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"fmt"
	"time"
)

var income = make(chan Document)

func TestWhether_correct_parameters_should_return_accepted(t *testing.T) {
	//in order to offer a reception to income
	go func() { <-income }()

	r := gofight.New()

	r.POST("/webhook").
	SetDebug(true).
	SetJSON(gofight.D{
		"_id": "1",
		"_rev": "1",
	}).
	Run(GinEngine(income), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, "{}\n", r.Body.String())
		assert.Equal(t, http.StatusAccepted, r.Code)
	})
}

func TestWhether_correct_parameters_should_propagate_to_income(t *testing.T) {

	go func() {
		message := Document{}
		message=<-income
		assert.Equal(t, "1", message.Id)
		assert.Equal(t, "1", message.Rev)
		fmt.Println("{end-of-test}")
	}()

	r := gofight.New()
	r.POST("/webhook").
	SetDebug(true).
	SetJSON(gofight.D{
		"_id": "1",
		"_rev": "1",
	}).
	Run(GinEngine(income), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

	})

	time.Sleep(100 * time.Millisecond)
}

func TestWhether_bad_parameters_should_response_bad_request(t *testing.T) {

	r := gofight.New()
	r.POST("/webhook").
	SetDebug(true).
	SetJSON(gofight.D{
		"no-param": "1",
	}).
	Run(GinEngine(income), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

