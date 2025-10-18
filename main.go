package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBookById).Methods("GET")

	r.HandleFunc("/book", createBook).Methods("POST")
	r.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")

	return r
}

func main() {
	// create router handler
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello World!")
	})

	router.HandleFunc("/books",getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBookById).Methods("GET")

	router.HandleFunc("/book", createBook).Methods("POST")
	router.HandleFunc("/book/{id}", updateBook).Methods("PUT")

	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")

	fmt.Print("Server is running on http://localhost:8080")

	// create connection and pass in the router handle (registered routes)
	http.ListenAndServe(":8080", router)
}