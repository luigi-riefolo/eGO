package main

import (
	"errors"
	"fmt"
)

//The most commonly-used error implementation is
// the errors package's unexported errorString type.

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

func errorFn() (int, error) {
	f := -1.0
	if f < 0 {
		return 0, fmt.Errorf("math: square root of negative number %g", f)
	}
	return 0, nil
}

func main() {
	_ = errors.New("math: square root of negative number")

	// Or
	//errorFn()
}

/*

 The error interface requires only a Error method; specific error implementations might have additional methods. For instance, the net package returns errors of type error, following the usual convention, but some of the error implementations have additional methods defined by the net.Error interface:

package net

type Error interface {
    error
    Timeout() bool   // Is the error a timeout?
    Temporary() bool // Is the error temporary?
}

Client code can test for a net.Error with a type assertion and then distinguish transient network errors from permanent ones. For instance, a web crawler might sleep and retry when it encounters a temporary error and give up otherwise.

if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
    time.Sleep(1e9)
    continue
}
if err != nil {
    log.Fatal(err)
}

*/
