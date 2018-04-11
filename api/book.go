package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Book Type
type Book struct {
	// define the book
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description,omitempty"`
}

func (b Book) ToJSON() []byte {
	ToJSON, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

func FromJSON(data []byte) Book {
	book := Book{}
	err := json.Unmarshal(data, &book)
	if err != nil {
		panic(err)
	}
	return book
}

/*
func writeJSON(w http.ResponseWriter, books map[string]Book) {
	b, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset-utf-8")
	w.Write(b)
}
*/
func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

// static Books Data for now
var books = map[string]Book{
	"1234567":   Book{Title: "Hitchhikers Guide", Author: "Douglas Adam", ISBN: "1234567", Description: "The best book ever"},
	"898989898": Book{Title: "Cloud Native Go", Author: "Some Guy", ISBN: "898989898"},
}

// static Books Data as array
var Books = []Book{
	Book{Title: "Hitchhikers Guide", Author: "Douglas Adam", ISBN: "1234567", Description: "The best book ever"},
	Book{Title: "Cloud Native Go", Author: "Some Guy", ISBN: "898989898"},
}

func BooksHandleFunc(w http.ResponseWriter, r *http.Request) {
	// implementing logic for /api/books
	switch method := r.Method; method {
	case http.MethodGet:
		fmt.Println("MethodGet (s)")
		books := AllBooks()
		writeJSON(w, books)
	case http.MethodPost:
		fmt.Println("MethodPost (s)")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		isbn, created := CreateBook(book)
		if created {
			w.Header().Add("Location", "/api/books/"+isbn)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		fmt.Println(method)
		fmt.Println("Default")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}
func BookHandleFunc(w http.ResponseWriter, r *http.Request) {
	// implementing logic for /api/books/<isbn>
	isbn := r.URL.Path[len("/api/book/"):]
	switch method := r.Method; method {
	case http.MethodGet:
		book, found := GetBook(isbn)
		if found {
			writeJSON(w, book)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		exists := UpdateBook(isbn, book)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		DeleteBook(isbn)
		w.WriteHeader(http.StatusOK)
	default:
		fmt.Println(method)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}

	b, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset-utf-8")
	w.Write(b)
}

// returns a slice of all books.
func AllBooks() map[string]Book {
	return books
}

// create a new book if it does not exist.
func CreateBook(book Book) (string, bool) {
	if val, ok := books[book.ISBN]; ok {
		return val.ISBN, false
	}
	books[book.ISBN] = book
	return book.ISBN, true
}

func GetBook(isbn string) (Book, bool) {
	if val, ok := books[isbn]; ok {
		return val, false
	}
	return books[isbn], true
}

// UpdateBook updates an existing book
func UpdateBook(isbn string, book Book) bool {
	_, exists := books[isbn]
	if exists {
		books[isbn] = book
	}
	return exists
}

// DeleteBook removes a book from the map by ISBN key
func DeleteBook(isbn string) {
	delete(books, isbn)
}
