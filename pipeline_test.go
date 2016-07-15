package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Channels struct {
	cdoc chan Document
	cerr chan error
}

func newContext(cdoc chan Document, cerr chan error) Channels {
	var out Channels
	out.cdoc = cdoc
	out.cerr = cerr
	return out
}


func sourceCanceled(in chan Document) Channels {
	cancel := make(chan error)
	close(cancel)
	//done <- "" // this provokes a panic in runtime
	return newContext(in, cancel)
}

func sourceOk() Channels {
	err := make(chan error)
	in := make(chan Document)
	return newContext(in, err)
}

func mapper(context Channels, handler func(Document) chan Document) Channels {
	select {
	case <-context.cerr:
		return context
	case doc:=<-context.cdoc:
		return newContext(handler(doc), context.cerr)
	}
}

func TestSourceCanceled(t *testing.T) {
	context := sourceCanceled(income)
	select {
	case <-context.cerr:
		t.Log("passed by case <-err:\n")
	case <-context.cdoc:
		t.Fail()
	}
}


func TestMapperOverCanceledSource(t *testing.T) {
	context := sourceCanceled(income)

	context2 := mapper(sourceCanceled(income), func(doc Document) chan Document {
		return context.cdoc
	})

	context3 := mapper(context2, func(doc Document) chan Document {
		return context2.cdoc
	})

	select {
	case <-context3.cerr:
		t.Log("passed by case <-err:\n")
	case <-context.cdoc:
		t.Fail()
	}

}

func TestMapperOverSourceOk(t *testing.T) {
	var doc Document
	doc.Id = "hi"
	context:=sourceOk()

	go func() { context.cdoc<-doc }()

	context2 := mapper(context, func(doc Document) chan Document {
		out:=make(chan Document)
		if doc.Id=="hi" {
			doc.Id="hi2"
		} else {
			doc.Id="bad"
		}
		go func() { out<-doc }()
		return out
	})



	context3 := mapper(context2, func(doc Document) chan Document {
		out:=make(chan Document)
		if doc.Id=="hi2" {
			doc.Id="hi3"

		} else {
			doc.Id="bad"
		}
		go func() { out<-doc }()
		return out
	})

	select {
	case <-context3.cerr:
		t.Fail()
	case result:=<-context3.cdoc:
		t.Log("doc mapped to "+result.Id)
		assert.Equal(t, "hi3", result.Id)
	}

}

func TestMapperOverSourceOkButCanceled(t *testing.T) {
	var doc Document
	doc.Id = "hi"
	context:=sourceOk()

	go func() { context.cdoc<-doc }()

	context2 := mapper(context, func(doc Document) chan Document {
		out:=make(chan Document)
		if doc.Id=="hi" {
			doc.Id="hi2"
		} else {
			doc.Id="bad"
		}
		go func() { out<-doc }()
		return out
	})



	context3 := mapper(context2, func(doc Document) chan Document {
		out:=make(chan Document)
		if doc.Id=="hi2" {
			close(context2.cerr)
		} else {
			doc.Id="bad"
		}
		return out
	})

	select {
	case <-context3.cerr:
		t.Log("passed by case <-err:\n")
	case <-context3.cdoc:
		t.Fail()
	}

}
