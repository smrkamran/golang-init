package main

import (
	"fmt"
	"net/http"
)

type DB interface {
	Store(string) error
}

type Store struct{}

func (s *Store) Store(value string) error {
	fmt.Println("Storing into the db...")
	return nil
}

func myExecuteFunc(db DB) ExecuteFn {
	return func(s string) {
		fmt.Println("my ex func, ", s)
		db.Store(s)
	}
}

func makeHttpFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.Store("Some Http thingy...")
	}
}

func main() {
	s := &Store{}

	http.HandleFunc("/", makeHttpFunc(s))

	Execute(myExecuteFunc(s))
}

// COming from Third party lib
type ExecuteFn func(string)

func Execute(fn ExecuteFn) {
	fn("FOO BAR BAZ")
}
