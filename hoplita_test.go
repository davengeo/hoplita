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

func first(in chan Document) (chan Document, chan string) {
	done := make(chan string)
	close(done)
	//done <- "" // this provokes a panic in runtime
	return in, done
}

type pepe func() chan Document


func second(in chan Document, done chan string, handler chan Document) (chan Document, chan string) {
	select {
		case <-done:
			return in, done
		case <-in:
			return handler, done
	}
}

func TestPipeLine(t *testing.T) {
	in, err := first(income)
	select {
		case <-err:
			t.Log("passed by case <-err:")
		case <-in:
			t.Fail()
	}
}


func TestSecondPipeLine(t *testing.T) {
	in, done := first(income)

	doc:=Document{}

	sec, err := second(in, done, func(doc Document) chan Document {
		out:=make(chan Document)
		return out
	}(doc))


	select {
		case <-err:
			t.Log("passed by case <-err:")
		case <-sec:
			t.Fail()
	}

}