package main

import "fmt"

// Stringer ...
type Stringer interface {
	String() string
}
type fakeString struct {
	content string
}

/*
if err != nil {
	if msqlerr, ok := err.(*mysql.MySQLError); ok && msqlerr.Number == 1062 {
		log.Println("We got a MySQL duplicate :(")
	}
}
*/

// function used to implement the Stringer interface
func (s *fakeString) String() string {
	return s.content
}
func printString(value interface{}) {
	switch str := value.(type) {
	case string:
		fmt.Println("1111")
		fmt.Println(str)
	case Stringer:
		fmt.Println("2222")
		fmt.Println(str.String())
	}
}
func main() {
	s := &fakeString{"Ceci n'est pas un string"}
	printString(s)
	printString("Hello, Gophers")
}
