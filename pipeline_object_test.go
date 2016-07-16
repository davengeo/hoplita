package main

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)


func newChannels(cdoc chan Document, cerr chan error) Channels {
	var out Channels
	out.cdoc = cdoc
	out.cerr = cerr
	return out
}

func (c Channels) mapper2(handler func(Document) Document) Channels {
	select {
	case <-c.cerr:
		return c
	case doc:=<-c.cdoc:
		out:=make(chan Document)
		go func() {
			newDoc:=handler(doc)
			out<-newDoc
		}()
		return newChannels(out, c.cerr)
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


func TestCreationChannelsHappy(t *testing.T) {
	err := make(chan error)
	in := make(chan Document)
	context := newChannels(in, err)

	var doc Document
	doc.Id = "hi"

	go func() { context.cdoc<-doc }()

	context.
		mapper2(func(doc Document) Document {
			doc.Id+="hi2"
			return doc
		}).
		mapper2(func(doc Document) Document {
			doc.Id+="hi3"
			return doc
		}).
		resolve(func(doc Document) {
			t.Log("passed with value:"+doc.Id)
			assert.Equal(t, "hihi2hi3", doc.Id)
		}, func(err error) {
			t.Fail()
		})

}

func TestCreationChannelsCanceled(t *testing.T) {
	err := make(chan error)
	in := make(chan Document)
	context := newChannels(in, err)

	var doc Document
	doc.Id = "hi"

	go func() { context.cdoc<-doc }()

	context.
	mapper2(func(doc Document) Document {
		doc.Id+="hi2"
		er := errors.New("error")
		context.cerr<-er
		//cancellation
		close(context.cerr)
		return doc
	}).
	mapper2(func(doc Document) Document {
		doc.Id+="hi3"
		return doc
	}).
	resolve(func(doc Document) {
		t.Log("passed with value:"+doc.Id)
		t.Fail()
	}, func(err error) {
		t.Log("passed with error")
	})

}