package main

import "testing"

type Context struct {
	cdoc chan Document
	cerr chan error
}

func newContext(cdoc chan Document, cerr chan error) *Context {
	out:= new(Context)
	out.cdoc = cdoc
	out.cerr = cerr
	return out
}

func first(in chan Document) *Context {
	done := make(chan error)
	close(done)
	//done <- "" // this provokes a panic in runtime
	return newContext(in, done)
}

func second(context *Context, handler chan Document) *Context {
	select {
	case <-context.cerr:
		return context
	case <-context.cdoc:
		return newContext(handler, context.cerr)
	}
}

func TestPipeLine(t *testing.T) {
	context := first(income)
	select {
	case <-context.cerr:
		t.Log("passed by case <-err:\n")
	case <-context.cdoc:
		t.Fail()
	}
}


func TestSecondPipeLine(t *testing.T) {
	context := first(income)

	context2 := second(context, func() chan Document {
		return context.cdoc
	}())

	context3 := second(context2, func() chan Document {
		return context2.cdoc
	}())

	select {
	case <-context3.cerr:
		t.Log("passed by case <-err:\n")
	case <-context.cdoc:
		t.Fail()
	}

}
