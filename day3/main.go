package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
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

type Error struct {
	Message string `json:"Error"`
}

type Info struct {
	Info string `json:"Info"`
}

type Item struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}

func main() {
	fmt.Println("API was started")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/items", GetAllItems).Methods("GET")
	myRouter.HandleFunc("/item/{id}", GetItemById).Methods("GET")

	myRouter.HandleFunc("/item", AddNewItem).Methods("POST")

	myRouter.HandleFunc("/item/{id}", UpdateItemById).Methods("PUT")

	myRouter.HandleFunc("/item/{id}", DeleteItemById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))

}

func GetAllItems(writer http.ResponseWriter, request *http.Request) {
	items := selectAll()
	if len(items) < 1 {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "No one items found in DB"})
	} else {
		json.NewEncoder(writer).Encode(items)
	}
}

func selectAll() []Item {
	db := ConnectToDb()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса %v", err)
	}
	defer rows.Close()

	var allItems []Item
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.Id, &item.Amount, &item.Price, &item.Title)
		if err != nil {
			fmt.Printf("Что-то пошло не так при чтении строки с результатом %v\r\n", err)
			continue
		}
		allItems = append(allItems, item)
	}

	return allItems
}

func GetItemById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "Couldn't convert string to int"})
		return
	}

	item := getItemFromDb(id)
	if item != nil {
		json.NewEncoder(writer).Encode(item)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "This item doesn't exist"})
	}
}

func getItemFromDb(id int) *Item {
	db := ConnectToDb()
	defer db.Close()

	var item Item
	err := db.QueryRow("SELECT * FROM items WHERE id = $1", id).Scan(&item.Id, &item.Amount, &item.Price, &item.Title)
	if err != nil {
		fmt.Printf("Возникла ошибка при выполнении запроса %v\r\n", err)
		return nil
	}
	return &item
}

func AddNewItem(writer http.ResponseWriter, request *http.Request) {
	var item Item
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &item)

	db := ConnectToDb()
	defer db.Close()

	var newId int
	err := db.QueryRow("INSERT INTO items(amount, price, title) VALUES ($1,$2,$3) RETURNING id", item.Amount, item.Price, item.Title).Scan(&newId)
	if err != nil {
		fmt.Printf("Возникла ошибка при выполнении запроса %v\r\n", err)
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "Не удалось внести элемент в справочник."})
	} else {
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(Info{Info: fmt.Sprintf("Успешно создан элемент с id = %v", newId)})
	}
}

func UpdateItemById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	var item Item
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &item)
	id, _ := strconv.Atoi(vars["id"])

	if updateItemById(id, item) {
		writer.WriteHeader(http.StatusAccepted)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func updateItemById(id int, item Item) bool {
	db := ConnectToDb()
	defer db.Close()

	rows, err := db.Exec("UPDATE items SET amount=$1, price=$2, title=$3 WHERE id=$4", item.Amount, item.Price, item.Title, id)
	if err != nil {
		fmt.Printf("Возникла ошибка при выполнении запроса %v\r\n", err)
		return false
	}
	if val, _ := rows.RowsAffected(); val == 0 {
		fmt.Println("Нет такого id")
		return false
	}
	return true
}

func DeleteItemById(writer http.ResponseWriter, request *http.Request) {
	idStr := mux.Vars(request)["id"]
	id, _ := strconv.Atoi(idStr)
	if deleteItemFromDb(id) {
		writer.WriteHeader(http.StatusAccepted)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func deleteItemFromDb(id int) bool {
	db := ConnectToDb()
	defer db.Close()
	res, err := db.Exec("DELETE FROM items WHERE id=$1", id)
	if err != nil {
		log.Fatalf("Ошибка даления записи по id %v", err)
	}

	if val, _ := res.RowsAffected(); val == 0 {
		fmt.Println("Нет такого id")
		return false
	}
	return true
}

func ConnectToDb() *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Что-то не подключилось %v", err)
	}
	return db
}
