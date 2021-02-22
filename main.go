package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"math/rand"

	//"math/rand"
	"net/http"
	"strconv"

	pg "github.com/go-pg/pg"
	orm "github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
)

//BookItem struct (MODEL)
type BookItem struct {
	tableName  struct{} `sql:"book_collection"`
	RefPointer int      `sql:"-"` //SQL will ignore this
	ID         int      `sql:"id,pk"`
	Isbn       int      `sql:"isbn,type:integer"`
	Title      string   `sql:"title,unique"`
	Image      string   `sql:"image"`
	Price      float64  `sql:"price,type:real"`
	IsActive   bool     `sql:"is_active,type:boolean"`
	Author     *Author  `sql:"author,type:json"`
}

//Author is selected
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//Response is
type Response struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

//Book is
type Book struct {
	ID     int     `json:"id"`
	Isbn   int     `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
	Price  float64 `json:"price"`
}

var gDb *pg.DB

func main() {
	log.Printf(("Starting"))

	//r.HandleFunc("/api/books", getBooks).Methods("GET")
	//r.HandleFunc("/api/books/{id}", getBook).Methods("GET")

	//r.HandleFunc("/api/books/{id}", updateBooks).Methods("PUT")
	//r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	//
	gDb = connect()
	if gDb != nil {

	}
	r := mux.NewRouter()
	r.HandleFunc("/api/books", createBook).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Scanln()
	disconnect(gDb)
	return
}

func connect() *pg.DB {
	options := &pg.Options{
		User:     "postgres",
		Password: "kamran92",
		Addr:     "localhost:5432",
	}
	db := pg.Connect(options)

	if db == nil {
		log.Printf("Failed to connect to database.\n")

	}
	log.Printf("Connected to database")

	createBooksTable(db)
	return db
}

func disconnect(db *pg.DB) {
	closeErr := db.Close()
	if closeErr != nil {
		panic("Error while closing the connection:")

	}
	log.Printf("Connection closed")
	return
}

func createBooksTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	createErr := db.CreateTable(&BookItem{}, opts)
	if createErr != nil {
		log.Printf("Error while creating table")
		return createErr
	}
	log.Printf("Table BookItems is created")
	return nil
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	//_ = json.NewDecoder(r.Body).Decode(&book)
	var author *Author
	var autString = r.FormValue("author")
	json.Unmarshal([]byte(autString), &author)
	book.Author = author
	book.Title = r.FormValue("title")
	//books = append(books, book)
	prc, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	id, _ := strconv.Atoi(r.FormValue("id"))
	book.Price = prc
	book.Isbn = isbn
	book.ID = id

	saveBook(gDb, &book)
	resp := Response{
		Result:  true,
		Message: "done",
	}
	json.NewEncoder(w).Encode(resp)
}

func saveBook(dbRef *pg.DB, book *Book) {

	newBook := BookItem{
		ID:       book.ID,
		Isbn:     book.Isbn,
		Title:    book.Title,
		Price:    book.Price,
		IsActive: true,
		Author:   book.Author,
	}
	insertErr := dbRef.Insert(newBook)
	if insertErr != nil {
		panic(insertErr)
	}
	log.Printf("Book inserted")
}
