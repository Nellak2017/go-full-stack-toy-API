package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
)

// Model for our Book API
type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book // make a slice of Book structs
var db *sql.DB

func init() {
	gotenv.Load() // loads the environment variables
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))

	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)

	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter() // Returns a new router instance

	router.HandleFunc("/books", getBooks).Methods("GET")           // When this endpoint is reached with a GET request, it will return the getBooks method
	router.HandleFunc("/books/{id}", getBook).Methods("GET")       // When this endpoint is reached with a GET request, it will return the getBook method
	router.HandleFunc("/books", addBook).Methods("POST")           // When this endpoint is reached with a POST request, it will return the addBook method
	router.HandleFunc("/books", updateBook).Methods("PUT")         // When this endpoint is reached with a PUT request, it will return the updateBook method
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE") // When this endpoint is reached with a DELETE request, it will return the removeBook method

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query("SELECT * FROM books") // Get all books in the database table
	logFatal(err)

	defer rows.Close() // Close db after all logic executed

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year) // ??? See documentation. My guess is that Scan modifies the book model by putting db items inside it, using the pointer stuff.
		logFatal(err)
		books = append(books, book) // add book to the list
	}

	json.NewEncoder(w).Encode(books) // When this endpoint is reached, write all the books in our library to the Response object
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r) // get our params from the request. Ex: map[id:1]. --> This is a map object containing a key of id and a String value 1. 1 must be converted to int to be used
	log.Println(params)

	row := db.QueryRow("SELECT * FROM books WHERE id=$1", params["id"])

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book) // then write to the Response the book we requested (Then write that book's value to our Response object)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book) // I find passing &book to be confusing. Why do I do this in this situation and not others? Pointers are retarded as fuck.

	err := db.QueryRow("INSERT INTO books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID) // Add a book to the DB

	logFatal(err)

	json.NewEncoder(w).Encode(bookID) // When this endpoint is reached, write the added ID to the Response object
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id;", &book.Title, &book.Author, &book.Year, &book.ID) // Update a book in the DB

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("DELETE FROM books WHERE id=$1;", params["id"]) // Add a book to the DB
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
