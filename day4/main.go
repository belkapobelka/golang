package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var config = Config{
	User:     "guest",
	Password: "guest",
	DBName:   "store",
	SSLMode:  "disable",
}
var connStr = fmt.Sprintf("user=%v password=%v database=%v sslmode=%v", config.User, config.Password, config.DBName, config.SSLMode)

type ErrorMessage struct {
	Message string `json:"message"`
}

type Item struct {
	ID     int    `json:"id"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
	Title  string `json:"title"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/items", GetAll).Methods("GET")

	router.HandleFunc("/item/{id}", DeleteItemById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func DeleteItemById(writer http.ResponseWriter, request *http.Request) {
	idStr := request.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if DeleteItemFromDb(id) {
		writer.WriteHeader(http.StatusAccepted)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func DeleteItemFromDb(id int) bool {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Что-то не подключилось %v", err)
	}
	defer db.Close()

	res, err := db.Exec(`DELETE FROM items WHERE id=$1`, id)
	if err != nil {
		log.Fatalf("Ошибка даления записи по id %v", err)
	}

	if val, _ := res.RowsAffected(); val == 0 {
		return false
	}
	return true
}

func GetAll(writer http.ResponseWriter, request *http.Request) {
	items := SelectAll()
	if len(items) < 1 {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(ErrorMessage{Message: "No one items found in DB"})
	} else {
		json.NewEncoder(writer).Encode(items)
	}
}

func SelectAll() []Item {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Что-то не подключилось %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса %v", err)
	}
	defer rows.Close()

	var allItems []Item
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Amount, &item.Price, &item.Title)
		if err != nil {
			fmt.Printf("Что-то пошло не так при чтении строки %v", err)
			continue
		}
		allItems = append(allItems, item)
	}

	return allItems
}
