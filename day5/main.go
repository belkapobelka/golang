package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	jwtmiddlware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var config = Config{
	Port:     ":8080",
	User:     "guest",
	Password: "guest",
	DBName:   "store",
	SSLMode:  "disable",
}
var connStr = fmt.Sprintf("user=%v password=%v database=%v sslmode=%v", config.User, config.Password, config.DBName, config.SSLMode)

type Error struct {
	Message string `json:"message"`
}

type Article struct {
	Id      int    `json:"Id"`
	Content string `json:"Content"`
}

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var Articles = []Article{{1, "First"}}

var SecretKey = []byte("secret") // == apiKey

func main() {
	fmt.Println("API was started")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/articles", GetAllArticles).Methods("GET")
	// use auth
	router.Handle("/article", jwtMiddleware.Handler(http.HandlerFunc(AddNewArticle))).Methods("POST")

	router.HandleFunc("/auth", PostToken).Methods("POST")
	router.HandleFunc("/register", RegisterUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)

	InsertUser(user)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration OK"))
}

func InsertUser(user User) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO Users (login,password) VALUES ($1,$2)", user.Login, user.Password)
	if err != nil {
		log.Fatal(err)
	}
}

func PostToken(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)

	Users := SelectAll()
	for _, u := range Users {
		if u.Login == user.Login && u.Password == user.Password {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims) // map with params

			claims["admin"] = true
			claims["name"] = "New User"
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

			tokenString, err := token.SignedString(SecretKey)
			if err != nil {
				log.Fatal(err)
			}
			w.Write([]byte(tokenString))
			return
		}
	}

	w.WriteHeader(401)
	w.Write([]byte("You are not in User Database"))
}

func SelectAll() []User {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		log.Fatal(err)
	}

	var Users []User
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Login, &u.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Users = append(Users, u)
	}
	return Users
}

var jwtMiddleware = jwtmiddlware.New(jwtmiddlware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func AddNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	Articles = append(Articles, article)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Articles)
}
