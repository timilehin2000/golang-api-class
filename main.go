package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type library struct {
	dbHost, dbPassword, dbName string
}

type Book struct {
	Id   string
	Name string
	ISBN string
}

const (
	ApiPath = "/api/v1/books"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "timmy419"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = ApiPath
	}

	l := library{
		dbHost:     dbHost,
		dbName:     dbName,
		dbPassword: dbPassword,
	}

	l.openConnection()

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.GetBooks).Methods("GET")

	fmt.Println("Server is running")
	http.ListenAndServe(":8080", r)
}

func (l library) GetBooks(w http.ResponseWriter, r *http.Request) {
	db := l.openConnection()

	rows, err := db.Query("SELECT * FROM books")

	if err != nil {
		log.Fatalf("Error loading books %s\n", err.Error())
	}

	books := []Book{}

	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)

		if err != nil {
			log.Fatalf("Scanning the rows %s\n", err.Error())
		}

		book := Book{
			Id:   id,
			Name: name,
			ISBN: isbn,
		}

		books = append(books, book)
		json.NewEncoder(w).Encode(books)

		l.closeConnection(db)
	}

}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", l.dbPassword, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("Cannot open DB connection %s\n", err.Error())
	}

	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()

	if err != nil {
		log.Fatalf("Closing DB %s/n", err.Error())
	}
}
