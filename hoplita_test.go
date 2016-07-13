package main

import (
	"testing"
	"gopkg.in/appleboy/gofight.v1"
	"github.com/stretchr/testify/assert"
	"net/http"
)


func TestWhether_correct_parameters_should_return_accepted(t *testing.T) {

	r := gofight.New()
	r.POST("/webhook").
	SetDebug(true).
	SetJSON(gofight.D{
		"_id": "1",
		"_rev": "1",
	}).
	Run(GinEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, "{}\n", r.Body.String())
		assert.Equal(t, http.StatusAccepted, r.Code)
	})

}