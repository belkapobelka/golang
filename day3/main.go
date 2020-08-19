package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Error struct {
	Message string `json:"Error"`
}

type Item struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}

var Items []Item

func main() {
	Items = []Item{
		{
			Id:     1,
			Item:   "Item 1",
			Amount: 11,
			Price:  "111.11",
		},
		{
			Id:     2,
			Item:   "Item 2",
			Amount: 22,
			Price:  "222.22",
		},
		{
			Id:     3,
			Item:   "Item 3",
			Amount: 33,
			Price:  "333.33",
		},
	}

	fmt.Println("API was started")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/items", getAllItems).Methods("GET")
	myRouter.HandleFunc("/item/{id}", getItemById).Methods("GET")

	myRouter.HandleFunc("/item", addNewItem).Methods("POST")

	myRouter.HandleFunc("/item/{id}", updateItemById).Methods("PUT")

	myRouter.HandleFunc("/item/{id}", deleteItemById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))

}

func getAllItems(writer http.ResponseWriter, request *http.Request) {
	if len(Items) != 0 {
		json.NewEncoder(writer).Encode(Items)
	} else {
		err := Error{Message: "No one item exists"}
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(err)
	}
}

func getItemById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "Couldn't convert string to int"})
		return
	}

	for _, item := range Items {
		if item.Id == id {
			json.NewEncoder(writer).Encode(item)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode(Error{Message: "This item doesn't exist"})
}

func addNewItem(writer http.ResponseWriter, request *http.Request) {
	var item Item
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &item)
	Items = append(Items, item)
	writer.WriteHeader(http.StatusCreated)
}

func updateItemById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	reqBody, _ := ioutil.ReadAll(request.Body)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "Couldn't convert string to int"})
		return
	}

	for index, item := range Items {
		if item.Id == id {
			json.Unmarshal(reqBody, &Items[index])
			writer.WriteHeader(http.StatusAccepted)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode(Error{Message: "Item with this id doesn't exist"})
}

func deleteItemById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(Error{Message: "Couldn't convert string to int"})
		return
	}

	for index, item := range Items {
		if item.Id == id {
			Items = append(Items[:index], Items[index+1:]...)
			writer.WriteHeader(http.StatusAccepted)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode(Error{Message: "Item with this id doesn't exist"})
}
