package main

import "testing"

func first(in chan Document) (chan Document, chan string) {
	done := make(chan string)
	close(done)
	//done <- "" // this provokes a panic in runtime
	return in, done
}

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
		t.Log("passed by case <-err:\n")
	case <-in:
		t.Fail()
	}
}


func TestSecondPipeLine(t *testing.T) {
	in, done := first(income)

	sec, er2 := second(in, done, func() chan Document {
		return in
	}())

	third, er3 := second(sec, er2, func() chan Document {
		return in
	}())

	select {
	case <-er3:
		t.Log("passed by case <-err:\n")
	case <-third:
		t.Fail()
	}

}
