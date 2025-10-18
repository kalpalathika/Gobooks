package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}

// in-memory mock data
var books = []Book{
    {ID: "1", Title: "Book1", Author: "Author1", Year: 1997},
    {ID: "2", Title: "Book2", Author: "Author2", Year: 1998},
    {ID: "3", Title: "Book3", Author: "Author3", Year: 1999},
	{ID: "4", Title: "Book4", Author: "Author4", Year: 1993},
	{ID: "5", Title: "Book5", Author: "Author5", Year: 1999},
	{ID: "6", Title: "Book6", Author: "Author6", Year: 1995},
	{ID: "7", Title: "Book7", Author: "Author7", Year: 1991},
	{ID: "8", Title: "Book8", Author: "Author8", Year: 1990},
}

type ApiResponse struct{
	Success bool `json:"success"`
	Data interface{} `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}


// Helpers 

func writeJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// headers are already sent; just log it
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ApiResponse{
		Success: false,
		Error:   msg,
	})
}

func validateBook (newBook Book) []string{
	var missingFields []string

	if newBook.Title == "" {
		missingFields = append(missingFields, "title")
	}
	if newBook.Author == "" {
		missingFields = append(missingFields, "author")
	}
	if newBook.Year <= 0 {
		missingFields = append(missingFields, "year")
	}

	return missingFields
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool{

	// 2MB cap to avoid huge bodies
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err:= dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}

	return true
}

// Main functions
func getBooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, ApiResponse{
		Success: true,
		Data:    books,
	})
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, book := range books {
		if book.ID == vars["id"] {
			writeJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data: book,
			})
			return
		}
	}

	writeError(w, http.StatusNotFound, "Book not found")
}

func createBook(w http.ResponseWriter, r *http.Request){
    var newBook Book
	if !decodeJSON(w, r, &newBook) {
		return
	}

	// --- Validation ---
	if missing := validateBook(newBook); len(missing) > 0 {
		writeError(w, http.StatusBadRequest,
			fmt.Sprintf("Missing or invalid fields: %v", strings.Join(missing, ", ")))
		return
	}

	// Temporary ID strategy
	newBook.ID = strconv.Itoa(len(books) + 1)
	books = append(books, newBook)

	// RESTful: 201 + Location + return the created resource
	w.Header().Set("Location", "/books/"+newBook.ID)
	writeJSON(w, http.StatusCreated, ApiResponse{
		Success: true,
		Data:    newBook,
	})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var updated Book
	if !decodeJSON(w, r, &updated) { // uses DisallowUnknownFields + 1MB cap + Close
		return
	}

	if missing := validateBook(updated); len(missing) > 0 {
		writeError(w, http.StatusBadRequest,
			fmt.Sprintf("Missing or invalid fields: %v", strings.Join(missing, ", ")))
		return
	}

	for i, b := range books {
		if b.ID == id {
			updated.ID = b.ID // path param 
			books[i] = updated
			writeJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    updated,
			})
			return
		}
	}

	writeError(w, http.StatusNotFound, "Book Id does not exist")
}


func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)

			// success: 204 No Content, no body
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeError(w, http.StatusNotFound, "Book Id does not exist")
}



