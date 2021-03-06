package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Author  string `json:"Author"`
	Content string `json:"Content"`
}

var Articles []Article

func main() {
	Articles = []Article{
		{
			Id:      "1",
			Title:   "First Title",
			Author:  "First Author",
			Content: "First Content",
		},
		{
			Id:      "2",
			Title:   "Second Title",
			Author:  "Second Author",
			Content: "Second Content",
		},
	}

	fmt.Println("API started")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/articles", getAllArticles).Methods("GET")
	myRouter.HandleFunc("/article/{id}", getArticleById).Methods("GET")

	myRouter.HandleFunc("/article", addNewArticle).Methods("POST")

	myRouter.HandleFunc("/article/{id}", deleteArticleById).Methods("DELETE")

	myRouter.HandleFunc("/article/{id}", updateArticleById).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func updateArticleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: updateArticleById worked ...")

	vars := mux.Vars(r)
	for index, article := range Articles {
		if article.Id == vars["id"] {
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &Articles[index])
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Articles[index])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	err := Error{Message: "Статья с указанным id не найдена."}
	json.NewEncoder(w).Encode(err)
}

func deleteArticleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: deleteArticleById worked ...")

	vars := mux.Vars(r)
	flag := false

	for index, article := range Articles {
		if article.Id == vars["id"] {
			flag = true
			Articles = append(Articles[:index], Articles[index+1:]...)
			w.WriteHeader(http.StatusAccepted)
		}
	}

	if !flag {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: "Статья с указанным id не найдена."}
		json.NewEncoder(w).Encode(err)
	}
}

func addNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: addNewArticle worked ...")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	Articles = append(Articles, article)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func getArticleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: getArticleById worked ...")
	vars := mux.Vars(r)
	find := false

	for _, article := range Articles {
		if article.Id == vars["id"] {
			find = true
			json.NewEncoder(w).Encode(article)
		}
	}

	if !find {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: "Статья с указанным id не найдена."}
		json.NewEncoder(w).Encode(err)
	}
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: getAllArticles worked ...")
	json.NewEncoder(w).Encode(Articles)
}
