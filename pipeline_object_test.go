package main

import (
	"testing"
	_ "github.com/stretchr/testify/assert"
	"fmt"
	"github.com/stretchr/testify/assert"
)


func newChannels(cdoc chan Document, cerr chan error) Channels {
	var out Channels
	out.cdoc = cdoc
	out.cerr = cerr
	return out
}

func (c Channels) mapper2(handler func(Document) chan Document) Channels {
	select {
	case <-c.cerr:
		return c
	case doc:=<-c.cdoc:
		return newChannels(handler(doc), c.cerr)
	}
}

func (c Channels) resolve(done func(Document), fail func(error)) {
	select {
	case err:=<-c.cerr:
		fail(err)
	case doc:=<-c.cdoc:
		done(doc)
	}
}


func TestCreationChannels(t *testing.T) {
	err := make(chan error)
	in := make(chan Document)
	context := newChannels(in, err)

	var doc Document
	doc.Id = "hi"

	go func() { context.cdoc<-doc }()

	context.
		mapper2(func(doc Document) chan Document {
			out:=make(chan Document)
			doc.Id+="hi2"
			go func() { out<-doc }()
			return out
		}).
		mapper2(func(doc Document) chan Document {
			out:=make(chan Document)
			doc.Id+="hi3"
			go func() { out<-doc }()
			return out
		}).
		resolve(func(doc Document) {
			t.Log("passed with value:"+doc.Id)
			assert.Equal(t, "hihi2hi3", doc.Id)
		}, func(err error) {
			fmt.Println(err)
			t.Fail()
		})

}